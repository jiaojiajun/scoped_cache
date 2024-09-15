package dependency_utils

import (
	"github.com/gin-gonic/gin"
	"reflect"
)

func SetDependency[T any](ctx *gin.Context, dependency *T) {
	if ctx == nil {
		panic(" context is nil")
	}
	var container = GetScopedContainer(ctx)
	if container == nil {
		panic("container is nil")
	}
	var dependencyKey = GenDependencyKey[T]()
	if container.dependencies[dependencyKey] != nil {
		panic("dependency already exists")
	}
	container.dependencies[dependencyKey] = dependency
}

func GenDependencyKey[T any]() string {
	var ptr *T
	var t = reflect.TypeOf(ptr).Elem()
	var name = t.Name()
	if name == "" {
		panic("dependency name is empty")
	}
	name = t.PkgPath() + "::" + name
	return name
}

func GetDependency[T any](ctx *gin.Context) *T {
	if ctx == nil {
		panic("ctx is nil")
	}
	var key = GenDependencyKey[T]()
	var container = GetScopedContainer(ctx)
	if container.dependencies[key] == nil {
		panic("dependency is not set")
	}
	return container.dependencies[key].(*T)
}
