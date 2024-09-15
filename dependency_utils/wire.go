package dependency_utils

import (
	"github.com/gin-gonic/gin"
	"github.com/jiaojiajun/scoped_cache_demo/service"
)

func WireApp(ctx *gin.Context) {
	ctx = WithScopedContainer(ctx)
	ctx = WireScopedContainer(ctx)
}

func WireScopedContainer(ctx *gin.Context) *gin.Context {
	var nameScopedCacheService = service.NewNameScopedCacheService()
	SetDependency[service.NameScopedCacheService](ctx, nameScopedCacheService)
	return ctx
}
