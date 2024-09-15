package cache

import (
	"github.com/gin-gonic/gin"
	"github.com/vmihailenco/msgpack"
	"golang.org/x/sync/singleflight"
	"golang.org/x/xerrors"
	"sync"
)

type CacheType int32

const (
	CacheType_ObjectCache    CacheType = 0
	CacheType_ByteArrayCache CacheType = 1
)

type ScopedCache[TKey any, TValue any] struct {
	cacheType              CacheType
	fnGenKey               func(TKey) string
	fnGetValueWithoutCache func(*gin.Context, TKey) (*TValue, error)

	singleflight singleflight.Group
	cacheLock    sync.Mutex
	objectCache  map[string]*TValue
	byteCache    map[string][]byte
}

func NewScopedCache[Tkey any, TValue any](
	cacheType CacheType,
	fnGenKey func(Tkey) string,
	fnGetValueWithoutCache func(ctx *gin.Context, key Tkey) (*TValue, error),
) *ScopedCache[Tkey, TValue] {
	var x = &ScopedCache[Tkey, TValue]{
		cacheType:              cacheType,
		fnGenKey:               fnGenKey,
		fnGetValueWithoutCache: fnGetValueWithoutCache,
	}
	switch cacheType {
	case CacheType_ObjectCache:
		x.objectCache = make(map[string]*TValue)
	case CacheType_ByteArrayCache:
		x.byteCache = make(map[string][]byte)
	default:
		panic("unknown cache type")
	}
	return x

}

func (p *ScopedCache[TKey, TValue]) readCache(key string) *TValue {
	p.cacheLock.Lock()
	defer p.cacheLock.Unlock()

	switch p.cacheType {
	case CacheType_ObjectCache:
		return p.objectCache[key]
	case CacheType_ByteArrayCache:
		var bytes = p.byteCache[key]
		if len(bytes) == 0 {
			return nil
		}
		var result TValue
		var err = msgpack.Unmarshal(bytes, result)
		if err != nil {
			panic("msgpack unmarshal err")
		}
		return &result
	default:
		panic("unreachable method")

	}

}

func (p *ScopedCache[TKey, TValue]) writeCache(key string, value *TValue) {
	p.cacheLock.Lock()
	defer p.cacheLock.Unlock()

	if value == nil {
		panic("write nil value")
	}

	switch p.cacheType {
	case CacheType_ObjectCache:
		p.objectCache[key] = value
	case CacheType_ByteArrayCache:
		var bytes, err = msgpack.Marshal(value)
		if err != nil {
			panic(err)
		}
		p.byteCache[key] = bytes
	default:
		panic("unreachable method")
	}

}

func (p *ScopedCache[TKey, TValue]) Fetch(ctx *gin.Context, tKey TKey) (*TValue, error) {
	var key = p.fnGenKey(tKey)
	resultObj, err, shared := p.singleflight.Do(key, func() (interface{}, error) {
		var result = p.readCache(key)
		if result != nil {
			return result, nil
		}
		result, err := p.fnGetValueWithoutCache(ctx, tKey)
		if err != nil {
			return nil, xerrors.Errorf(": %w", err)
		}
		if result == nil {
			return nil, nil
		}
		p.writeCache(key, result)
		return result, nil
	})
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	var resultValue = resultObj.(*TValue)
	if shared && p.cacheType == CacheType_ByteArrayCache {
		resultValue = p.ByteCopy(resultValue)
	}
	return resultValue, nil
}

func (p *ScopedCache[TKey, TValue]) ByteCopy(value *TValue) *TValue {
	if value == nil {
		return nil
	}
	bytes, err := msgpack.Marshal(value)
	if err != nil {
		panic(err)
	}
	var res TValue
	err = msgpack.Unmarshal(bytes, res)
	if err != nil {
		panic(err)
	}
	return &res
}
