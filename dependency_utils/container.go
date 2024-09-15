package dependency_utils

import "github.com/gin-gonic/gin"

type Container struct {
	dependencies map[string]any
}

func newContainer() *Container {
	return &Container{dependencies: map[string]any{}}
}

const defaultContainerKey = "default"

func WithScopedContainer(ctx *gin.Context) *gin.Context {
	var _, exists = ctx.Get(defaultContainerKey)
	if exists {
		panic("duplicate container key")
	}
	ctx.Set(defaultContainerKey, newContainer())
	return ctx
}

func GetScopedContainer(ctx *gin.Context) *Container {
	var value, exists = ctx.Get(defaultContainerKey)
	if !exists || value == nil {
		panic("container not exists")
	}
	var container = value.(*Container)
	if container == nil {
		panic("container is nil")
	}
	return container
}
