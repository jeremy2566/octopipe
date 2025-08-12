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
