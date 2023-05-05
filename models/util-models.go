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
	ProjectId        string                  `json:"projectId,omitempty" yaml:"projectId" bson:"projectId,omitempty"`
	AppName          string                  `json:"appName,omitempty" yaml:"appName" bson:"appName,omitempty"`
	Description      string                  `json:"description,omitempty" yaml:"description" bson:"description,omitempty"`
	ImageURL         string                  `json:"imageUrl,omitempty" yaml:"imageURL" bson:"imageURL,omitempty"`
	Version          string                  `json:"version,omitempty" yaml:"version" bson:"version,omitempty"`
	HostName         string                  `json:"hostName,omitempty" yaml:"hostName" bson:"hostName,omitempty"`
	ClusterURL       string                  `json:"clusterURL,omitempty" yaml:"clusterURL" bson:"clusterURL,omitempty"`
	ClusterIPs       []string                `json:"clusterIPs,omitempty" yaml:"clusterIPs" bson:"clusterIPs,omitempty"`
	Microservices    map[string]Microservice `json:"microservices,omitempty" yaml:"microservices" bson:"microservices,omitempty"`
	Monitoring       bool                    `json:"monitoring,omitempty" yaml:"monitoring" bson:"monitoring,omitempty"`
	EvaluationResult EvaluationResponse      `json:"evaluationResult,omitempty" yaml:"evaluationResult" bson:"evaluationResult,omitempty"`
}

type Microservice struct {
	ServiceName       string                `json:"serviceName,omitempty" yaml:"serviceName" bson:"serviceName,omitempty"`
	Configs           []string              `json:"configs,omitempty" yaml:"configs" bson:"configs,omitempty"`
	AvgReplicas       int                   `json:"avgReplicas,omitempty" yaml:"avgReplicas" bson:"avgReplicas,omitempty"`
	MinReplicas       int                   `json:"minReplicas,omitempty" yaml:"minReplicas" bson:"minReplicas,omitempty"`
	MaxReplicas       int                   `json:"maxReplicas,omitempty" yaml:"maxReplicas" bson:"maxReplicas,omitempty"`
	MaxCPU            string                `json:"maxCPU,omitempty" yaml:"maxCPU" bson:"maxCPU,omitempty"`
	MaxMemory         string                `json:"maxMemory,omitempty" yaml:"maxMemory" bson:"maxMemory,omitempty"`
	DockerImage       string                `json:"dockerImage,omitempty" yaml:"dockerImage" bson:"dockerImage,omitempty"`
	ContainerPort     int                   `json:"containerPort,omitempty" yaml:"containerPort" bson:"containerPort,omitempty"`
	Envs              map[string]EnvRequest `json:"envs,omitempty" yaml:"envs" bson:"envs,omitempty"`
	ServiceEvaluation ServiceEvaluation     `json:"serviceEvaluation,omitempty" yaml:"serviceEvaluation" bson:"serviceEvaluation,omitempty"`
}

type EvaluationResponse struct {
	Language                string   `json:"language,omitempty" bson:"language,omitempty"`
	Database                string   `json:"database,omitempty" bson:"database,omitempty"`
	HasDockerized           bool     `json:"hasDockerized" bson:"hasDockerized"`
	HasKubernetesService    bool     `json:"hasKubernetesService" bson:"hasKubernetesService"`
	HasKubernetesDeployment bool     `json:"hasKubernetesDeployment" bson:"hasKubernetesDeployment"`
	Microservices           []string `json:"microservices,omitempty" bson:"microservices,omitempty"`
}

type EnvRequest struct {
	Name  string `json:"name,omitempty" yaml:"name" bson:"name,omitempty"`
	Value string `json:"value,omitempty" yaml:"value" bson:"value,omitempty"`
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
