package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a Api) ViewAllNamespace(c *gin.Context) {
	namespaces := a.cache.ViewAll()
	c.JSON(http.StatusOK, namespaces)
}

func (a Api) DeleteNamespace(c *gin.Context) {
	subEnv := c.Param("sub_env")
	err := a.cache.DeleteNamespace(subEnv)
	print(err)
}
