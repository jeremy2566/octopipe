package dao

import (
	"testing"

	"github.com/jeremy2566/octopipe/internal/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestRedisImpl_SaveNamespace(t *testing.T) {
	log, _ := zap.NewDevelopment()
	r := NewRdb(log)

	err := r.SaveNamespace("test-1", model.DaoNamespace{
		SubEnv:      "test-1",
		Branch:      "branch-1",
		ServiceName: []string{"service-1"},
	})
	assert.Nil(t, err)
}

func TestNewRedis_GetNamespace(t *testing.T) {
	log, _ := zap.NewDevelopment()
	r := NewRdb(log)
	namespace, _ := r.GetValueByKey("test-1")
	log.Info("namespace", zap.Any("namespace", namespace))
}

func TestNewRedis_GetAllNamespace(t *testing.T) {
	log, _ := zap.NewDevelopment()
	r := NewRdb(log)
	namespace := r.GetAllNamespace()
	log.Info("namespace", zap.Any("ns", namespace))
}

func TestNewRdb_DeleteNamespace(t *testing.T) {
	log, _ := zap.NewDevelopment()
	r := NewRdb(log)
	_ = r.DeleteNamespace("test-1")
}

func TestNewRdb_GetNamespaceByBranch(t *testing.T) {
	log, _ := zap.NewDevelopment()
	rdb := NewRdb(log)
	namespace, _ := rdb.GetNamespaceByBranch("branch1-1")
	assert.NotNil(t, namespace)
	//log.Info("TestNewRdb_GetNamespaceByBranch", zap.Any("namespace", namespace))
}
