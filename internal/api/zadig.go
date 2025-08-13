package api

import (
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
	err := a.zadig.CreateSubEnv()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "create sub env failed.",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "create sub env success.",
	})
}

func (a Api) ServiceCharts(c *gin.Context) {
	charts := a.zadig.GetServiceCharts()
	c.JSON(http.StatusOK, gin.H{
		"data": charts,
	})
}
