package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jiaojiajun/scoped_cache_demo/cache"
	"golang.org/x/xerrors"
	"time"
)

type Name struct {
	FistName   string
	FamilyName string
}
type NameScopedCacheService struct {
	nameCache *cache.ScopedCache[string, Name]
}

func NewNameScopedCacheService() *NameScopedCacheService {
	var p = &NameScopedCacheService{}
	p.initCache()
	return p
}
func (p *NameScopedCacheService) initCache() {

	nameCache := cache.NewScopedCache[string, Name](
		cache.CacheType_ObjectCache,
		func(i string) string {
			return fmt.Sprint(i)
		},
		p.GetName_NoCache,
	)

	p.nameCache = nameCache
}

var singletonNameScopedCacheService = initSingletonNameScopedCacheService()

func initSingletonNameScopedCacheService() *NameScopedCacheService {
	return NewNameScopedCacheService()
}

func GetSingletonNameScopedCacheService() *NameScopedCacheService {
	return singletonNameScopedCacheService
}

func (p *NameScopedCacheService) GetName_NoCache(ctx *gin.Context, key string) (*Name, error) {
	//log.Print("called no cache method \n")
	time.Sleep(time.Second * 2)
	_ = fmt.Sprint(key) // 这一行如果不加会导致nocache方法不会被调用
	return &Name{
		FistName:   "name",
		FamilyName: "jiao",
	}, nil
}

func (p *NameScopedCacheService) GetNameWithCache(ctx *gin.Context, key string) (*Name, error) {
	name, err := p.nameCache.Fetch(ctx, key)
	if err != nil {
		return nil, xerrors.Errorf(":%w", err)
	}
	return name, nil
}
