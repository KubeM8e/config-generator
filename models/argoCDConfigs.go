package models

type ArgoCDApplicationConfig struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
	} `yaml:"metadata"`
	Spec struct {
		Project string `yaml:"project"`
		Source  struct {
			RepoURL        string `yaml:"repoURL"`
			TargetRevision string `yaml:"targetRevision"`
			Path           string `yaml:"path"`
			Helm           struct {
				ValueFiles []string `yaml:"valueFiles"`
			} `yaml:"helm"`
		} `yaml:"source"`
		Destination struct {
			Server    string `yaml:"server"`
			Namespace string `yaml:"namespace"`
		} `yaml:"destination"`
		SyncPolicy struct {
			SyncOptions []string `yaml:"syncOptions"`
			Automated   struct {
				SelfHeal bool `yaml:"selfHeal"`
				Prune    bool `yaml:"prune"`
			} `yaml:"automated"`
		} `yaml:"syncPolicy"`
	} `yaml:"spec"`
}
