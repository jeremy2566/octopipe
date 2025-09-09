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

type Callback struct {
	ObjectKind string `json:"object_kind"`
	Event      string `json:"event"`
	Workflow   struct {
		TaskID              int    `json:"task_id"`
		ProjectName         string `json:"project_name"`
		ProjectDisplayName  string `json:"project_display_name"`
		WorkflowName        string `json:"workflow_name"`
		WorkflowDisplayName string `json:"workflow_display_name"`
		Status              string `json:"status"`
		Remark              string `json:"remark"`
		DetailURL           string `json:"detail_url"`
		Error               string `json:"error"`
		CreateTime          int    `json:"create_time"`
		StartTime           int    `json:"start_time"`
		EndTime             int    `json:"end_time"`
		Stages              []struct {
			Name      string `json:"name"`
			Status    string `json:"status"`
			StartTime int    `json:"start_time"`
			EndTime   int    `json:"end_time"`
			Jobs      []struct {
				Name        string `json:"name"`
				DisplayName string `json:"display_name"`
				Type        string `json:"type"`
				Status      string `json:"status"`
				StartTime   int    `json:"start_time"`
				EndTime     int    `json:"end_time"`
				Error       string `json:"error"`
				Spec        struct {
					Repositories []struct {
						Source        string      `json:"source"`
						RepoOwner     string      `json:"repo_owner"`
						RepoNamespace string      `json:"repo_namespace"`
						RepoName      string      `json:"repo_name"`
						Branch        string      `json:"branch"`
						Prs           interface{} `json:"prs"`
						Tag           string      `json:"tag"`
						CommitID      string      `json:"commit_id"`
						CommitURL     string      `json:"commit_url"`
						CommitMessage string      `json:"commit_message"`
					} `json:"repositories"`
					Image string `json:"image"`
				} `json:"spec"`
			} `json:"jobs"`
			Error string `json:"error"`
		} `json:"stages"`
		TaskCreator      string `json:"task_creator"`
		TaskCreatorID    string `json:"task_creator_id"`
		TaskCreatorPhone string `json:"task_creator_phone"`
		TaskCreatorEmail string `json:"task_creator_email"`
		TaskType         string `json:"task_type"`
	} `json:"workflow"`
}

type RespTaskDetail struct {
	TaskID       int    `json:"task_id"`
	WorkflowKey  string `json:"workflow_key"`
	WorkflowName string `json:"workflow_name"`
	Params       []struct {
		Name         string      `json:"name"`
		Description  string      `json:"description"`
		Type         string      `json:"type"`
		Value        string      `json:"value"`
		Repo         interface{} `json:"repo"`
		ChoiceOption []string    `json:"choice_option"`
		Default      string      `json:"default"`
		IsCredential bool        `json:"is_credential"`
	} `json:"params"`
	Status      string `json:"status"`
	Reverted    bool   `json:"reverted"`
	Remark      string `json:"remark"`
	TaskCreator string `json:"task_creator"`
	TaskRevoker string `json:"task_revoker"`
	CreateTime  int    `json:"create_time"`
	StartTime   int    `json:"start_time"`
	EndTime     int    `json:"end_time"`
	Stages      []struct {
		Name       string      `json:"name"`
		Status     string      `json:"status"`
		StartTime  int         `json:"start_time"`
		EndTime    int         `json:"end_time"`
		Parallel   bool        `json:"parallel"`
		ManualExec interface{} `json:"manual_exec"`
		Jobs       []struct {
			Name             string `json:"name"`
			Key              string `json:"key"`
			DisplayName      string `json:"display_name"`
			OriginName       string `json:"origin_name"`
			Type             string `json:"type"`
			Status           string `json:"status"`
			Reverted         bool   `json:"reverted"`
			StartTime        int    `json:"start_time"`
			EndTime          int    `json:"end_time"`
			CostSeconds      int    `json:"cost_seconds"`
			Error            string `json:"error"`
			BreakpointBefore bool   `json:"breakpoint_before"`
			BreakpointAfter  bool   `json:"breakpoint_after"`
			Spec             struct {
				Repos []struct {
					Source        string `json:"source"`
					RepoOwner     string `json:"repo_owner"`
					RepoNamespace string `json:"repo_namespace"`
					RepoName      string `json:"repo_name"`
					RemoteName    string `json:"remote_name"`
					Branch        string `json:"branch"`
					EnableCommit  bool   `json:"enable_commit"`
					CommitID      string `json:"commit_id"`
					CommitMessage string `json:"commit_message"`
					Hidden        bool   `json:"hidden"`
					IsPrimary     bool   `json:"is_primary"`
					CodehostID    int    `json:"codehost_id"`
					OauthToken    string `json:"oauth_token"`
					Address       string `json:"address"`
					FilterRegexp  string `json:"filter_regexp"`
					SourceFrom    string `json:"source_from"`
					ParamName     string `json:"param_name"`
					JobName       string `json:"job_name"`
					ServiceName   string `json:"service_name"`
					ServiceModule string `json:"service_module"`
					RepoIndex     int    `json:"repo_index"`
					SubmissionID  string `json:"submission_id"`
					DisableSsl    bool   `json:"disable_ssl"`
				} `json:"repos"`
				Image         string `json:"image"`
				Package       string `json:"package"`
				ServiceName   string `json:"service_name"`
				ServiceModule string `json:"service_module"`
				Envs          []struct {
					Key          string `json:"key"`
					Value        string `json:"value"`
					Type         string `json:"type"`
					RegistryID   string `json:"registry_id"`
					IsCredential bool   `json:"is_credential"`
					Description  string `json:"description"`
				} `json:"envs"`
			} `json:"spec"`
			ErrorPolicy          interface{} `json:"error_policy"`
			ErrorHandlerUserID   string      `json:"error_handler_user_id"`
			ErrorHandlerUsername string      `json:"error_handler_username"`
			RetryCount           int         `json:"retry_count"`
			JobInfo              struct {
				JobName       string `json:"job_name"`
				ServiceModule string `json:"service_module"`
				ServiceName   string `json:"service_name"`
			} `json:"job_info"`
		} `json:"jobs"`
		Error string `json:"error"`
	} `json:"stages"`
	ProjectKey       string `json:"project_key"`
	IsRestart        bool   `json:"is_restart"`
	Debug            bool   `json:"debug"`
	ApprovalTicketID string `json:"approval_ticket_id"`
	ApprovalID       string `json:"approval_id"`
}
