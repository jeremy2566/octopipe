package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jeremy2566/octopipe/internal/model"
)

func (a Api) Allocator(c *gin.Context) {
	var req model.AllocatorReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := a.zadig.Allocator(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "success",
	})
}

func (a Api) Callback(c *gin.Context) {

}

func (a Api) CreateSubEnv(c *gin.Context) {
	subEnv, err := a.zadig.CreateSubEnv()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "create sub env failed.",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": fmt.Sprintf("create sub env[%s] success.", subEnv),
	})
}

func (a Api) AddService(c *gin.Context) {
	var req model.AddServiceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := a.zadig.AddService(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "success",
	})
}
func (a Api) DeployService(c *gin.Context) {
	var req model.DeployServiceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if _, err := a.zadig.DeployService(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "success",
	})
}

func (a Api) ServiceCharts(c *gin.Context) {
	charts := a.zadig.GetServiceCharts()
	c.JSON(http.StatusOK, gin.H{
		"data": charts,
	})
}

func (a Api) DeleteSubEnv(c *gin.Context) {
	subEnv := c.Param("sub_env")
	err := a.zadig.DeleteSubEnv(subEnv)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "success",
	})
}

func (a Api) Webhook(c *gin.Context) {
	var cb model.Callback
	if err := c.ShouldBindJSON(&cb); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := a.zadig.Webhook(cb)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "success",
	})
}
