package handler

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jeremy2566/octopipe/internal/cache"
	"go.uber.org/zap"
)

type NamespaceAllocatorReq struct {
	ServiceName string `json:"service_name" binding:"required"`
	BranchName  string `json:"branch_name" binding:"required"`
	GithubActor string `json:"github_actor" binding:"required"`
}

func (h Handler) NamespaceAllocator(c *gin.Context) {
	var req NamespaceAllocatorReq

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("bind json failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("received request for namespace allocation", zap.Any("params", req))

	switch {
	case req.BranchName == "main" || req.BranchName == "master":
		h.log.Info("master/main handle.")
	case req.BranchName == "release":
		h.log.Info("release handle.")
	case req.BranchName == "dev" || req.BranchName == "develop":
		h.log.Info("dev/develop handle.")
	case strings.HasPrefix(req.BranchName, "feat/") || strings.HasPrefix(req.BranchName, "feature/"):
		h.handleFeature(req)
	default:
		h.log.Warn("unknown branch name", zap.String("branch", req.BranchName))
	}
	return
}

func (h Handler) handleFeature(req NamespaceAllocatorReq) {
	c := cache.GetInstance()
	bn := req.BranchName

	h.log.Info("handle feature", zap.String("service_name", req.ServiceName), zap.String("branch_name", bn))
	if ns, b := c.Get(bn); b {
		h.log.Info("cache hit, deploying to namespace", zap.String("branch_name", bn), zap.String("ns", ns))
		// 增加服务
		err := h.AddServices(ns, req.ServiceName, req.BranchName)
		if err != nil {
			h.log.Error("add services failed",
				zap.String("service", req.ServiceName),
				zap.String("branch", req.BranchName),
				zap.Error(err),
			)
			h.Notify(fmt.Sprintf("namespace[%s] service[%s] add failed. msg: %s", ns, req.ServiceName, err.Error()))
			return
		}
		// 部署服务
		taskId, err := h.DeployService(ns, req.ServiceName, req.BranchName)
		if err != nil {
			h.log.Error("deploy services failed",
				zap.String("service", req.ServiceName),
				zap.String("branch", req.BranchName),
				zap.Error(err),
			)
			h.Notify(fmt.Sprintf("namespace[%s] service[%s] deploy failed.", ns, req.ServiceName))
		}
		h.Notify(fmt.Sprintf("namespace[%s] service[%s] deploy succeeded. task id = %d", ns, req.ServiceName, taskId))
	} else {
		h.log.Info("cache miss", zap.String("branch_name", bn))
		// 挑选一个 ns
		ns := h.SelectedNamespace()
		instance := cache.GetInstance()
		instance.Set(bn, ns)
		h.log.Info("selected ns.", zap.String("ns", ns))
		// 创建命名空间
		err := h.CreateSubEnv(ns, "redis-backoffice", "redis-general", "backoffice-v1-web-app", "bo-v1-assets", req.ServiceName)
		if err != nil {
			h.log.Error("create sub env failed", zap.Error(err))
			h.Notify(fmt.Sprintf("namespace[%s] created failed.", ns))
			return
		}

		// 初次部署服务
		_, err = h.DeployService(ns, "backoffice-v1-web", "master")
		if err != nil {
			h.log.Error("deploy backoffice-v1-web services failed",
				zap.Error(err),
			)
			h.Notify(fmt.Sprintf("namespace[%s] backoffice-v1-web-app deploy failed.", ns))
			return
		}
		taskId, err := h.DeployService(ns, req.ServiceName, req.BranchName)
		if err != nil {
			h.log.Error("deploy services failed",
				zap.String("service", req.ServiceName),
				zap.String("branch", req.BranchName),
				zap.Error(err),
			)
			h.Notify(fmt.Sprintf("namespace[%s] %s failed.", ns, req.BranchName))
			return
		}

		h.Notify(fmt.Sprintf("namespace[%s] created succeeded. task id = %d", ns, taskId))
	}
}

func (h Handler) SelectedNamespace() string {
	envs := h.GetZadigNamespace("fat-base-envrionment")

	usedNamespaces := make(map[string]bool)

	for _, e := range envs {
		usedNamespaces[e.Namespace] = true
	}
	h.log.Info("fetch environments success", zap.Any("usedNamespaces", usedNamespaces))

	// 存储所有可用的命名空间
	var availableNamespaces []string
	for i := 1; i <= 50; i++ {
		ns := fmt.Sprintf("test%d", i)
		if !usedNamespaces[ns] {
			availableNamespaces = append(availableNamespaces, ns)
		}
	}

	// 检查是否有可用的命名空间
	if len(availableNamespaces) == 0 {
		h.log.Warn("No available namespaces in the 1-50 pool.")
	}
	randomIndex := rand.Intn(len(availableNamespaces))
	return availableNamespaces[randomIndex]
}

func (h Handler) ViewCache(c *gin.Context) {
	instance := cache.GetInstance()
	items := instance.Items()

	c.JSON(http.StatusOK, items)
}
