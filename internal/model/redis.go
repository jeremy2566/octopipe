package model

type DaoNamespace struct {
	SubEnv      string   `json:"sub_env"`
	UpdateBy    string   `json:"update_by"`
	Branch      string   `json:"branch"`
	ServiceName []string `json:"service_name"`
}
