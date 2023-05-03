package configs

import (
	"config-generator/pkg/utils"
	"github.com/go-git/go-git/v5"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

func CreatePrometheusStackValuesYaml(clusterIPs []string, repoName string) (*git.Worktree, *git.Repository) {
	gitWorkTree, gitRepo := utils.CloneGitHubRepo(repoName, TmpArgoFolder)

	values := make(map[string]interface{})
	file, err := os.ReadFile("prometheusStack/values.yaml")
	if err != nil {
		log.Printf("Could not read promethues stack values.yaml: %s", err)
	}
	err = yaml.Unmarshal(file, values)
	if err != nil {
		log.Printf("Could not unmarshal promethues stack values.yaml: %s", err)
	}

	for _, v := range values {
		if m, ok := v.(map[interface{}]interface{}); ok {
			if _, okE := m["endpoints"]; okE {
				m["endpoints"] = clusterIPs
			}
		}
	}

	valuesYamlFile, err := os.Create(TmpArgoFolder + "/monitoringValues.yaml")
	if err != nil {
		log.Printf("Could not create argo dir: %s", err)
	}

	valuesYamlData, err := yaml.Marshal(&values)
	if err != nil {
		log.Printf("Could not marshal argo : %s", err)
	}

	_, err = valuesYamlFile.Write(valuesYamlData)
	if err != nil {
		log.Printf("Could not write argo : %s", err)
	}

	utils.PushToGitHub(gitWorkTree, gitRepo, []string{"monitoringValues.yaml"})

	return gitWorkTree, gitRepo
}
