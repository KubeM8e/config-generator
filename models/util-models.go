package models

// Manifests

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
		Replicas string `yaml:"replicas"`
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
	Name  string    `yaml:"name"`
	Image string    `yaml:"image"`
	Ports []PortDep `yaml:"ports"`
	Env   []Env     `yaml:"env"`
}

type PortDep struct {
	ContainerPort string `yaml:"containerPort"`
}

type Env struct {
	Name      string `yaml:"name"`
	ValueFrom struct {
		ConfigMapKeyRef struct {
			Name string `yaml:"name"`
			Key  string `yaml:"key"`
		} `yaml:"configMapKeyRef"`
		//SecretKeyRef struct {
		//	Name string `yaml:"name"`
		//	Key  string `yaml:"key"`
		//} `yaml:"secretKeyRef"`
	} `yaml:"valueFrom"`
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
		MinReplicas string    `yaml:"minReplicas"`
		MaxReplicas string    `yaml:"maxReplicas"`
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
	Name       string `yaml:"name"`
	Port       int    `yaml:"port"`
	TargetPort string `yaml:"targetPort"`
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
	AppName       string         `json:"appName" yaml:"appName"`
	Version       string         `json:"version" yaml:"version"`
	HostName      string         `json:"hostName" yaml:"hostName"`
	ClusterURL    string         `json:"clusterURL" yaml:"clusterURL"`
	Microservices []Microservice `json:"microservices" yaml:"microservices"`
}

type Microservice struct {
	ServiceName   string `json:"serviceName" yaml:"serviceName"`
	AvgReplicas   int    `json:"avgReplicas" yaml:"avgReplicas"`
	MinReplicas   int    `json:"minReplicas" yaml:"minReplicas"`
	MaxReplicas   int    `json:"maxReplicas" yaml:"maxReplicas"`
	DockerImage   string `json:"dockerImage" yaml:"dockerImage"`
	ContainerPort int    `json:"containerPort" yaml:"containerPort"`
	Envs          []struct {
		Name  string `json:"name" yaml:"name"`
		Value string `json:"value" yaml:"value"`
	} `json:"envs" yaml:"envs"`
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
