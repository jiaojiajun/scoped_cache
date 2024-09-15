package service_test

import (
	"github.com/bytedance/mockey"
	"github.com/gin-gonic/gin"
	"github.com/jiaojiajun/scoped_cache_demo/service"
	"github.com/jiaojiajun/scoped_cache_demo/utils"
	"testing"
)

func TestNameScopedCacheService_GetName_WithCache_Call_NoCache_Method_Once_when_called_twice(t *testing.T) {
	var m = mockey.Mock((*service.NameScopedCacheService).GetName_NoCache).Build()
	defer m.UnPatch()

	var assert = utils.ProdAssert(t)
	var ctx = &gin.Context{}

	// act
	var nameScopedCacheService = service.GetSingletonNameScopedCacheService()
	res, err := nameScopedCacheService.GetNameWithCache(ctx, "123")
	assert.NoError(err)
	assert.NotNil(res)
	res, err = nameScopedCacheService.GetNameWithCache(ctx, "123")
	assert.NoError(err)
	assert.NotNil(res)
	assert.Equal(m.Times(), 1)
}
