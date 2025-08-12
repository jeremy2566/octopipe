package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jeremy2566/octopipe/internal/api"
	"github.com/jeremy2566/octopipe/internal/cache"
	"github.com/jeremy2566/octopipe/internal/handler"
	"go.uber.org/zap"
)

func New(log *zap.Logger) *gin.Engine {
	r := gin.Default()
	h := handler.New(log)
	router := api.New(log)
	envs := h.GetZadigNamespace("fat-base-envrionment")
	instance := cache.GetInstance()
	for _, e := range envs {
		instance.Set(e.EnvKey, e.Namespace)
	}

	r.POST("/namespace_allocator", h.NamespaceAllocator)
	{
		r.GET("/cache/view", router.ViewAllNamespace)
		r.DELETE("/cache/:sub_env", router.DeleteNamespace)
		//r.PUT("/cache/sync", h.SyncCache)
	}
	r.POST("/zadig/webhook", h.Webhook)
	return r
}
