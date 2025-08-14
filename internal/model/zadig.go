package model

type RespZadigEnv struct {
	EnvKey     string `json:"env_key"`
	ClusterID  string `json:"cluster_id"`
	Namespace  string `json:"namespace"`
	Production bool   `json:"production"`
	RegistryID string `json:"registry_id"`
	Status     string `json:"status"`
	UpdateBy   string `json:"update_by"`
	UpdateTime int64  `json:"update_time"`
}

func (r RespZadigEnv) HasTest17OrTest33() bool {
	return r.Namespace == "test17" || r.Namespace == "test33"
}

type RespZadigEnvDetail struct {
	ClusterID       string        `json:"cluster_id"`
	EnvKey          string        `json:"env_key"`
	GlobalVariables []interface{} `json:"global_variables"`
	Namespace       string        `json:"namespace"`
	ProjectKey      string        `json:"project_key"`
	RegistryID      string        `json:"registry_id"`
	Services        []struct {
		Containers []struct {
			Image     string `json:"image"`
			ImagePath struct {
				Image string `json:"image"`
				Tag   string `json:"tag"`
			} `json:"imagePath"`
			ImageName string `json:"image_name"`
			Name      string `json:"name"`
			Type      string `json:"type"`
		} `json:"containers"`
		ServiceName string        `json:"service_name"`
		Status      string        `json:"status"`
		Type        string        `json:"type"`
		VariableKvs []interface{} `json:"variable_kvs"`
	} `json:"services"`
	Status     string `json:"status"`
	UpdateBy   string `json:"update_by"`
	UpdateTime int    `json:"update_time"`
}

func (r RespZadigEnvDetail) GetServices() []string {
	var ret []string
	for _, service := range r.Services {
		ret = append(ret, service.ServiceName)
	}
	return ret
}

type AllocatorReq struct {
	ServiceName string `json:"service_name" binding:"required"`
	BranchName  string `json:"branch_name" binding:"required"`
	GithubActor string `json:"github_actor" binding:"required"`
}

type DeployServiceReq struct {
	SubEnv      string `json:"sub_env" binding:"required"`
	ServiceName string `json:"service_name" binding:"required"`
	BranchName  string `json:"branch_name" binding:"required"`
	GithubActor string `json:"github_actor" binding:"required"`
}

type AddServiceReq struct {
	SubEnv      string `json:"sub_env" binding:"required"`
	ServiceName string `json:"service_name" binding:"required"`
}

type ChartInfoReq struct {
	EnvName        string `json:"envName"`
	ServiceName    string `json:"serviceName"`
	ChartVersion   string `json:"chartVersion"`
	DeployStrategy string `json:"deploy_strategy"`
}

type RespChartInfo struct {
	ServiceName  string `json:"service_name"`
	ChartVersion string `json:"chart_version"`
}

type ShareEnvReq struct {
	Enable  bool   `json:"enable"`
	IsBase  bool   `json:"isBase"`
	BaseEnv string `json:"base_env"`
}

type CreateSubEnvReq []struct {
	EnvName     string         `json:"env_name"`
	ClusterID   string         `json:"cluster_id"`
	RegistryID  string         `json:"registry_id"`
	ChartValues []ChartInfoReq `json:"chartValues"`
	Namespace   string         `json:"namespace"`
	IsExisted   bool           `json:"is_existed"`
	ShareEnv    ShareEnvReq    `json:"share_env"`
}

type DeployServicesReq struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Params      []struct {
		Name  string `json:"name"`
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"params"`
	Stages []struct {
		Name string `json:"name"`
		Jobs []struct {
			Name string `json:"name"`
			Type string `json:"type"`
			Spec struct {
				DefaultServiceAndBuilds []struct {
					ServiceName   string `json:"service_name"`
					ServiceModule string `json:"service_module"`
					KeyVals       []struct {
						Key   string `json:"key"`
						Value string `json:"value"`
						Type  string `json:"type"`
					} `json:"key_vals"`
					Repos []struct {
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
					KeyVals       []struct {
						Key   string `json:"key"`
						Value string `json:"value"`
						Type  string `json:"type"`
					} `json:"key_vals"`
					Repos []struct {
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
	} `json:"stages"`
	Project string `json:"project"`
}

type UtilsFun struct {
	ReplacePolicy string                `json:"replacePolicy"`
	EnvNames      []string              `json:"envNames"`
	ChartValues   []UtilsFunChartValues `json:"chartValues"`
}

type UtilsFunChartValues struct {
	EnvName         string `json:"envName"`
	ServiceName     string `json:"serviceName"`
	ReleaseName     string `json:"releaseName"`
	ChartVersion    string `json:"chartVersion"`
	Deploy_strategy string `json:"deploy_strategy"`
}
