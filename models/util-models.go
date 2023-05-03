package models

// Manifests

type Chart struct {
	APIVersion  string `yaml:"apiVersion""`
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
}

type ConfigMap struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Data map[string]string `yaml:"data"`
}

type Secret struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Type string            `yaml:"type"`
	Data map[string]string `yaml:"data"`
}

type Deployment struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name   string `yaml:"name"`
		Labels struct {
			App string `yaml:"app"`
		} `yaml:"labels"`
	} `yaml:"metadata"`
	Spec struct {
		Replicas int `yaml:"replicas"`
		Selector struct {
			MatchLabels struct {
				App string `yaml:"app"`
			} `yaml:"matchLabels"`
		} `yaml:"selector"`
		Template struct {
			Metadata struct {
				Labels struct {
					App string `yaml:"app"`
				} `yaml:"labels"`
			} `yaml:"metadata"`
			Spec struct {
				Containers []Container `yaml:"containers"`
			} `yaml:"spec"`
		} `yaml:"template"`
	} `yaml:"spec"`
}

type Container struct {
	Name         string         `yaml:"name"`
	Image        string         `yaml:"image"`
	Ports        []PortDep      `yaml:"ports"`
	Env          []Env          `yaml:"env"`
	VolumeMounts []VolumeMounts `yaml:"volumeMounts,omitempty"`
	Resources    struct {
		Requests struct {
			CPU    string `yaml:"cpu,omitempty"`
			Memory string `yaml:"memory,omitempty"`
		} `yaml:"requests,omitempty"`
		Limits struct {
			CPU    string `yaml:"cpu,omitempty"`
			Memory string `yaml:"memory,omitempty"`
		} `yaml:"limits,omitempty"`
	} `yaml:"resources,omitempty"`
}

type PortDep struct {
	ContainerPort int    `yaml:"containerPort"`
	Name          string `yaml:"name,omitempty"`
}

type Env struct {
	Name      string `yaml:"name"`
	ValueFrom struct {
		ConfigMapKeyRef struct {
			Name string `yaml:"name"`
			Key  string `yaml:"key"`
		} `yaml:"configMapKeyRef,omitempty"`
		SecretKeyRef struct {
			Name string `yaml:"name"`
			Key  string `yaml:"key"`
		} `yaml:"secretKeyRef,omitempty"`
	} `yaml:"valueFrom"`
}

type VolumeMounts struct {
	Name      string `yaml:"name"`
	MountPath string `yaml:"mountPath"`
}

type HorizontalPodAutoscaler struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		ScaleTargetRef struct {
			APIVersion string `yaml:"apiVersion"`
			Kind       string `yaml:"kind"`
			Name       string `yaml:"name"`
		} `yaml:"scaleTargetRef"`
		MinReplicas int       `yaml:"minReplicas"`
		MaxReplicas int       `yaml:"maxReplicas"`
		Metrics     []Metrics `yaml:"metrics"`
	} `yaml:"spec"`
}

type Metrics struct {
	Type     string `yaml:"type"`
	Resource struct {
		Name   string `yaml:"name"`
		Target struct {
			Type               string `yaml:"type"`
			AverageUtilization int    `yaml:"averageUtilization"`
		} `yaml:"target"`
	} `yaml:"resource"`
}

type VerticalPodAutoscaler struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		TargetRef struct {
			APIVersion string `yaml:"apiVersion"`
			Kind       string `yaml:"kind"`
			Name       string `yaml:"name"`
		} `yaml:"targetRef"`
		UpdatePolicy struct {
			UpdateMode string `yaml:"updateMode"`
		} `yaml:"updatePolicy"`
	} `yaml:"spec"`
}

type Service struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		Selector struct {
			App string `yaml:"app"`
		} `yaml:"selector"`
		Ports []PortSVC `yaml:"ports"`
		Type  string    `yaml:"type"`
	} `yaml:"spec"`
}

type PortSVC struct {
	Name       string `yaml:"name,omitempty"`
	Protocol   string `yaml:"protocol,omitempty"`
	Port       int    `yaml:"port"`
	TargetPort int    `yaml:"targetPort"`
}

type Ingress struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name        string `yaml:"name"`
		Annotations struct {
			RewriteTarget string `yaml:"nginx.ingress.kubernetes.io/rewrite-target"`
		} `yaml:"annotations"`
	} `yaml:"metadata"`
	Spec struct {
		Rules []Rule `yaml:"rules"`
	} `yaml:"spec"`
}

type Rule struct {
	Host string `yaml:"host"`
	HTTP struct {
		Paths []Path `yaml:"paths"`
	} `yaml:"http"`
}

type Path struct {
	Path     string `yaml:"path"`
	PathType string `yaml:"pathType"`
	Backend  struct {
		Service struct {
			Name string `yaml:"name"`
			Port struct {
				Name string `yaml:"name"`
			} `yaml:"port"`
		} `yaml:"service"`
	} `yaml:"backend"`
}

type PersistentVolume struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name   string `yaml:"name"`
		Labels struct {
			Type string `yaml:"type"`
		} `yaml:"labels"`
	} `yaml:"metadata"`
	Spec struct {
		StorageClassName string `yaml:"storageClassName"`
		Capacity         struct {
			Storage string `yaml:"storage"`
		} `yaml:"capacity"`
		AccessModes []string `yaml:"accessModes"`
		HostPath    struct {
			Path string `yaml:"path"`
		} `yaml:"hostPath"`
	} `yaml:"spec"`
}

type PersistentVolumeClaim struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name   string `yaml:"name"`
		Labels struct {
			App string `yaml:"app"`
		} `yaml:"labels"`
	} `yaml:"metadata"`
	Spec struct {
		AccessModes []string `yaml:"accessModes"`
		Resources   struct {
			Requests struct {
				Storage string `yaml:"storage"`
			} `yaml:"requests"`
		} `yaml:"resources"`
		Selector struct {
			MatchLabels struct {
				Type string `yaml:"type"`
			} `yaml:"matchLabels"`
		} `yaml:"selector"`
	} `yaml:"spec"`
}

type StatefulSet struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		ServiceName string `yaml:"serviceName"`
		Selector    struct {
			MatchLabels struct {
				App string `yaml:"app"`
			} `yaml:"matchLabels"`
		} `yaml:"selector"`
		Replicas int `yaml:"replicas"`
		Template struct {
			Metadata struct {
				Labels struct {
					App string `yaml:"app"`
				} `yaml:"labels"`
			} `yaml:"metadata"`
			Spec struct {
				Containers []Container `yaml:"containers"`
				Volumes    []Volume    `yaml:"volumes"`
			} `yaml:"spec"`
		} `yaml:"template"`
		VolumeClaimTemplates []VolumeClaimTemplate `yaml:"volumeClaimTemplates"`
	} `yaml:"spec"`
}

type Volume struct {
	Name                  string `yaml:"name"`
	PersistentVolumeClaim struct {
		ClaimName string `yaml:"claimName"`
	} `yaml:"persistentVolumeClaim"`
}

type VolumeClaimTemplate struct {
	Metadata struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		AccessModes []string `yaml:"accessModes"`
		Resources   struct {
			Requests struct {
				Storage string `yaml:"storage"`
			} `yaml:"requests"`
		} `yaml:"resources"`
		StorageClassName string `yaml:"storageClassName"`
	} `yaml:"spec"`
}

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
	ServiceName   string                `json:"serviceName,omitempty" yaml:"serviceName" bson:"serviceName,omitempty"`
	Configs       []string              `json:"configs,omitempty" yaml:"configs" bson:"configs,omitempty"`
	AvgReplicas   int                   `json:"avgReplicas,omitempty" yaml:"avgReplicas" bson:"avgReplicas,omitempty"`
	MinReplicas   int                   `json:"minReplicas,omitempty" yaml:"minReplicas" bson:"minReplicas,omitempty"`
	MaxReplicas   int                   `json:"maxReplicas,omitempty" yaml:"maxReplicas" bson:"maxReplicas,omitempty"`
	MaxCPU        string                `json:"maxCPU,omitempty" yaml:"maxCPU" bson:"maxCPU,omitempty"`
	MaxMemory     string                `json:"maxMemory,omitempty" yaml:"maxMemory" bson:"maxMemory,omitempty"`
	DockerImage   string                `json:"dockerImage,omitempty" yaml:"dockerImage" bson:"dockerImage,omitempty"`
	ContainerPort int                   `json:"containerPort,omitempty" yaml:"containerPort" bson:"containerPort,omitempty"`
	Envs          map[string]EnvRequest `json:"envs,omitempty" yaml:"envs" bson:"envs,omitempty"`
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
	Name  string `json:"name" yaml:"name"`
	Value string `json:"value" yaml:"value"`
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
