package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jiaojiajun/scoped_cache_demo/dependency_utils"
	"github.com/jiaojiajun/scoped_cache_demo/service"
	"golang.org/x/xerrors"
)

type NameHandler struct {
}

func NewNameHandler() *NameHandler {
	return &NameHandler{}
}
func initNameHandler() *NameHandler {
	return NewNameHandler()
}

var singletonNameHandler = initNameHandler()

func GetSingletonNameHandler() *NameHandler {
	return singletonNameHandler
}

func (p *NameHandler) GetFirstName(ctx *gin.Context) (string, error) {
	name, err := dependency_utils.GetDependency[service.NameScopedCacheService](ctx).GetNameWithCache(ctx, "id")
	if err != nil {
		return "", xerrors.Errorf(": %w", err)
	}
	return name.FistName, nil
}
