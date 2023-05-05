package models

// API

type CDPipelineRequest struct {
	GitHubRepoURL        string `json:"gitHubRepoURL"`
	TargetRevision       string `json:"targetRevision"`
	Path                 string `json:"path"`
	ApplicationName      string `json:"applicationName"`
	ApplicationNamespace string `json:"applicationNamespace"`
	ClusterURL           string `json:"clusterURL"`
}

type ConfigurationRequest struct {
	ProjectId        string                  `json:"projectId,omitempty" yaml:"projectId,omitempty" bson:"projectId,omitempty"`
	AppName          string                  `json:"appName,omitempty" yaml:"appName,omitempty" bson:"appName,omitempty"`
	Description      string                  `json:"description,omitempty" yaml:"description,omitempty" bson:"description,omitempty"`
	ImageURL         string                  `json:"imageUrl,omitempty" yaml:"imageURL,omitempty" bson:"imageURL,omitempty"`
	Version          string                  `json:"version,omitempty" yaml:"version,omitempty" bson:"version,omitempty"`
	HostName         string                  `json:"hostName,omitempty" yaml:"hostName,omitempty" bson:"hostName,omitempty"`
	ClusterURL       string                  `json:"clusterURL,omitempty" yaml:"clusterURL,omitempty" bson:"clusterURL,omitempty"`
	ClusterIPs       []string                `json:"clusterIPs,omitempty" yaml:"clusterIPs,omitempty" bson:"clusterIPs,omitempty"`
	Monitoring       bool                    `json:"monitoring,omitempty" yaml:"monitoring,omitempty" bson:"monitoring,omitempty"`
	Microservices    map[string]Microservice `json:"microservices,omitempty" yaml:"microservices,omitempty" bson:"microservices,omitempty"`
	EvaluationResult EvaluationResponse      `json:"evaluationResult,omitempty" yaml:"evaluationResult,omitempty" bson:"evaluationResult,omitempty"`
}

type Microservice struct {
	ServiceName       string                `json:"serviceName,omitempty" yaml:"serviceName,omitempty" bson:"serviceName,omitempty"`
	Configs           []string              `json:"configs,omitempty" yaml:"configs,omitempty" bson:"configs,omitempty"`
	AvgReplicas       int                   `json:"avgReplicas,omitempty" yaml:"avgReplicas,omitempty" bson:"avgReplicas,omitempty"`
	MinReplicas       int                   `json:"minReplicas,omitempty" yaml:"minReplicas,omitempty" bson:"minReplicas,omitempty"`
	MaxReplicas       int                   `json:"maxReplicas,omitempty" yaml:"maxReplicas,omitempty" bson:"maxReplicas,omitempty"`
	MaxCPU            string                `json:"maxCPU,omitempty" yaml:"maxCPU,omitempty" bson:"maxCPU,omitempty"`
	MaxMemory         string                `json:"maxMemory,omitempty" yaml:"maxMemory,omitempty" bson:"maxMemory,omitempty"`
	DockerImage       string                `json:"dockerImage,omitempty" yaml:"dockerImage,omitempty" bson:"dockerImage,omitempty"`
	ContainerPort     int                   `json:"containerPort,omitempty" yaml:"containerPort,omitempty" bson:"containerPort,omitempty"`
	Envs              map[string]EnvRequest `json:"envs,omitempty" yaml:"envs,omitempty" bson:"envs,omitempty"`
	ServiceEvaluation ServiceEvaluation     `json:"serviceEvaluation,omitempty" yaml:"serviceEvaluation,omitempty" bson:"serviceEvaluation,omitempty"`
}

type EvaluationResponse struct {
	Language                string   `json:"language,omitempty" yaml:"language,omitempty" bson:"language,omitempty"`
	Database                string   `json:"database,omitempty" yaml:"database,omitempty" bson:"database,omitempty"`
	HasDockerized           bool     `json:"hasDockerized" yaml:"hasDockerized,omitempty" bson:"hasDockerized"`
	HasKubernetesService    bool     `json:"hasKubernetesService" yaml:"hasKubernetesService,omitempty" bson:"hasKubernetesService"`
	HasKubernetesDeployment bool     `json:"hasKubernetesDeployment" yaml:"hasKubernetesDeployment,omitempty" bson:"hasKubernetesDeployment"`
	Microservices           []string `json:"microservices,omitempty" yaml:"microservices,omitempty" bson:"microservices,omitempty"`
}

type EnvRequest struct {
	Name  string `json:"name,omitempty" yaml:"name,omitempty" bson:"name,omitempty"`
	Value string `json:"value,omitempty" yaml:"value,omitempty" bson:"value,omitempty"`
}

type ServiceEvaluation struct {
	KubeConfigType []string `json:"kubeConfigType,omitempty" yaml:"kubeConfigType" bson:"kubeConfigType,omitempty"`
}

type RepoResponse struct {
	HelmRepo string `json:"helmRepo"`
	ArgoRepo string `json:"argoRepo"`
}

// Database

type CreateDBData struct {
	AppName string    `json:"appName" bson:"appName"`
	Apps    DBAppData `json:"apps" bson:"apps"`
}

type DBAppData struct {
	Version string `json:"version" bson:"version"`
	AppID   string `json:"appID" bson:"appID"`
	//Deployments []DeploymentData `json:"deployments"`
}

type DeploymentData struct {
	DeploymentName string `json:"deploymentName"`
	DeploymentID   string `json:"deploymentID"`
}

// GitHub

type GitHubRepoData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
	AutoInit    bool   `json:"auto_init"`
}

type KubeComponent struct {
	KubeComponentType string      `json:"kubeComponentType"`
	APIVersion        string      `json:"apiVersion"`
	Kind              string      `json:"kind"`
	KubeObjectValue   interface{} `json:"kubeObjectValue"`
	KubeObjectKey     string      `json:"kubeObjectKey"`
}
