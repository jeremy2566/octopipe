package handler

import (
	"encoding/json"
	"fmt"
	"testing"

	"go.uber.org/zap"
	"resty.dev/v3"
)

func TestName(t *testing.T) {
	namespace := "test4"
	serviceName := "auth-api"
	branchName := "WB-9297"
	req := DeployServicesReq{
		Name:        "test33",
		DisplayName: "fat-base-workflow",
		Project:     "fat-base-envrionment",
		Params: []struct {
			Name  string `json:"name"`
			Type  string `json:"type"`
			Value string `json:"value"`
		}{
			{
				Name:  "环境",
				Type:  "choice",
				Value: namespace,
			},
		},
		Stages: []struct {
			Name string `json:"name"`
			Jobs []struct {
				Name string `json:"name"`
				Type string `json:"type"`
				Spec struct {
					DefaultServiceAndBuilds []struct {
						ServiceName   string `json:"service_name"`
						ServiceModule string `json:"service_module"`
						Repos         []struct {
							Source        string `json:"source"`
							RepoOwner     string `json:"repo_owner"`
							RepoNamespace string `json:"repo_namespace"`
							RepoName      string `json:"repo_name"`
							RemoteName    string `json:"remote_name"`
							Branch        string `json:"branch"`
							CodehostID    int    `json:"codehost_id"`
						} `json:"repos"`
					} `json:"default_service_and_builds"`
					ServiceAndBuilds []struct {
						ServiceName   string `json:"service_name"`
						ServiceModule string `json:"service_module"`
						BuildName     string `json:"build_name"`
						Repos         []struct {
							Source        string `json:"source"`
							RepoOwner     string `json:"repo_owner"`
							RepoNamespace string `json:"repo_namespace"`
							RepoName      string `json:"repo_name"`
							RemoteName    string `json:"remote_name"`
							Branch        string `json:"branch"`
							CodehostID    int    `json:"codehost_id"`
						} `json:"repos"`
					} `json:"service_and_builds"`
				} `json:"spec"`
			} `json:"jobs"`
		}{
			{
				Name: "构建",
				Jobs: []struct {
					Name string `json:"name"`
					Type string `json:"type"`
					Spec struct {
						DefaultServiceAndBuilds []struct {
							ServiceName   string `json:"service_name"`
							ServiceModule string `json:"service_module"`
							Repos         []struct {
								Source        string `json:"source"`
								RepoOwner     string `json:"repo_owner"`
								RepoNamespace string `json:"repo_namespace"`
								RepoName      string `json:"repo_name"`
								RemoteName    string `json:"remote_name"`
								Branch        string `json:"branch"`
								CodehostID    int    `json:"codehost_id"`
							} `json:"repos"`
						} `json:"default_service_and_builds"`
						ServiceAndBuilds []struct {
							ServiceName   string `json:"service_name"`
							ServiceModule string `json:"service_module"`
							BuildName     string `json:"build_name"`
							Repos         []struct {
								Source        string `json:"source"`
								RepoOwner     string `json:"repo_owner"`
								RepoNamespace string `json:"repo_namespace"`
								RepoName      string `json:"repo_name"`
								RemoteName    string `json:"remote_name"`
								Branch        string `json:"branch"`
								CodehostID    int    `json:"codehost_id"`
							} `json:"repos"`
						} `json:"service_and_builds"`
					} `json:"spec"`
				}{
					{
						Name: "构建发布",
						Type: "zadig-build",
						Spec: struct {
							DefaultServiceAndBuilds []struct {
								ServiceName   string `json:"service_name"`
								ServiceModule string `json:"service_module"`
								Repos         []struct {
									Source        string `json:"source"`
									RepoOwner     string `json:"repo_owner"`
									RepoNamespace string `json:"repo_namespace"`
									RepoName      string `json:"repo_name"`
									RemoteName    string `json:"remote_name"`
									Branch        string `json:"branch"`
									CodehostID    int    `json:"codehost_id"`
								} `json:"repos"`
							} `json:"default_service_and_builds"`
							ServiceAndBuilds []struct {
								ServiceName   string `json:"service_name"`
								ServiceModule string `json:"service_module"`
								BuildName     string `json:"build_name"`
								Repos         []struct {
									Source        string `json:"source"`
									RepoOwner     string `json:"repo_owner"`
									RepoNamespace string `json:"repo_namespace"`
									RepoName      string `json:"repo_name"`
									RemoteName    string `json:"remote_name"`
									Branch        string `json:"branch"`
									CodehostID    int    `json:"codehost_id"`
								} `json:"repos"`
							} `json:"service_and_builds"`
						}{
							DefaultServiceAndBuilds: []struct {
								ServiceName   string `json:"service_name"`
								ServiceModule string `json:"service_module"`
								Repos         []struct {
									Source        string `json:"source"`
									RepoOwner     string `json:"repo_owner"`
									RepoNamespace string `json:"repo_namespace"`
									RepoName      string `json:"repo_name"`
									RemoteName    string `json:"remote_name"`
									Branch        string `json:"branch"`
									CodehostID    int    `json:"codehost_id"`
								} `json:"repos"`
							}{
								{
									ServiceName:   serviceName,
									ServiceModule: serviceName,
									Repos: []struct {
										Source        string `json:"source"`
										RepoOwner     string `json:"repo_owner"`
										RepoNamespace string `json:"repo_namespace"`
										RepoName      string `json:"repo_name"`
										RemoteName    string `json:"remote_name"`
										Branch        string `json:"branch"`
										CodehostID    int    `json:"codehost_id"`
									}{
										{
											Source:        "github",
											RepoOwner:     "storehubnet",
											RepoNamespace: "storehubnet",
											RepoName:      serviceName,
											RemoteName:    "origin",
											Branch:        branchName,
											CodehostID:    6,
										},
									},
								},
							},
							ServiceAndBuilds: []struct {
								ServiceName   string `json:"service_name"`
								ServiceModule string `json:"service_module"`
								BuildName     string `json:"build_name"`
								Repos         []struct {
									Source        string `json:"source"`
									RepoOwner     string `json:"repo_owner"`
									RepoNamespace string `json:"repo_namespace"`
									RepoName      string `json:"repo_name"`
									RemoteName    string `json:"remote_name"`
									Branch        string `json:"branch"`
									CodehostID    int    `json:"codehost_id"`
								} `json:"repos"`
							}{
								{
									ServiceName:   serviceName,
									ServiceModule: serviceName,
									BuildName:     fmt.Sprintf("fat-base-envrionment-build-%s-1", serviceName),
									Repos: []struct {
										Source        string `json:"source"`
										RepoOwner     string `json:"repo_owner"`
										RepoNamespace string `json:"repo_namespace"`
										RepoName      string `json:"repo_name"`
										RemoteName    string `json:"remote_name"`
										Branch        string `json:"branch"`
										CodehostID    int    `json:"codehost_id"`
									}{
										{
											Source:        "github",
											RepoOwner:     "storehubnet",
											RepoNamespace: "storehubnet",
											RepoName:      serviceName,
											RemoteName:    "origin",
											Branch:        branchName,
											CodehostID:    6,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	marshal, _ := json.Marshal(req)
	fmt.Println(string(marshal))
}

func TestHandler_DeployService(t *testing.T) {
	h := Handler{
		log:    zap.NewNop(),
		client: resty.New(),
	}
	taskId, err := h.DeployService("test15", "auth-api", "master")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(taskId)
}

func TestZadig_GetTaskDetail(t *testing.T) {

}
