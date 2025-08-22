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
		r.GET("/cache/view/:sub_env", router.ViewNamespace)
		r.DELETE("/cache/:sub_env", router.DeleteNamespace)
	}

	{
		r.POST("/github/allocator", router.Allocator)
	}

	{
		r.POST("/zadig/namespace", router.CreateSubEnv)
		r.DELETE("/zadig/namespace/:sub_env", router.DeleteSubEnv)
		r.POST("/zadig/add_service", router.AddService)
		r.POST("/zadig/deploy_service", router.DeployService)
		r.GET("/zadig/service_charts", router.ServiceCharts)
		r.POST("/zadig/webhook", router.Webhook)
	}
	return r
}
