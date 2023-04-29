package core

import (
	"config-generator/models"
	"config-generator/pkg/utils"
	"fmt"
	"github.com/go-git/go-git/v5"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	deploymentName  = "deployment"
	serviceName     = "service"
	ingressName     = "ingress"
	templatesFolder = "templates"
	tmpHelmFolder   = "helm"
)

var templatesFolderPath = tmpHelmFolder + "/" + templatesFolder

func GenerateValuesYamlFile(configs map[string]interface{}, repoName string) (*git.Worktree, *git.Repository) {
	gitWorkTree, gitRepo := utils.CloneGitHubRepo(repoName, tmpHelmFolder)

	// create a helm folder
	err := os.MkdirAll(tmpHelmFolder, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	//creates values.yaml file inside the helm folder
	valuesFile, err := os.Create(tmpHelmFolder + "/values.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// converts map to yaml file
	yamlData, _ := yaml.Marshal(&configs)
	_, err = valuesFile.Write(yamlData)
	if err != nil {
		log.Fatal(err)
	}

	return gitWorkTree, gitRepo
}

func ConfigureHelmChart(configs map[string]interface{}, gitWorkTree *git.Worktree, gitRepo *git.Repository) {

	// temporarily hold the response map in the generatePlaceholders function
	responseMap := make(map[string]interface{})

	// creates templates folder inside helm folder
	err := os.MkdirAll(templatesFolderPath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	for key, value := range configs {

		if strings.EqualFold(key, deploymentName) {
			deploymentObject := models.KubeComponent{
				KubeComponentType: "deployment",
				APIVersion:        "apps/v1",
				Kind:              "Deployment",
				KubeObjectValue:   value,
				KubeObjectKey:     key,
			}
			generateHelmTemplate(deploymentObject, responseMap)

		} else if strings.EqualFold(key, serviceName) {
			serviceObject := models.KubeComponent{
				KubeComponentType: "service",
				APIVersion:        "v1",
				Kind:              "Service",
				KubeObjectValue:   value,
				KubeObjectKey:     key,
			}
			generateHelmTemplate(serviceObject, responseMap)

		} else if strings.EqualFold(key, ingressName) {
			ingressObject := models.KubeComponent{
				KubeComponentType: "ingress",
				APIVersion:        "networking.k8s.io/v1",
				Kind:              "Ingress",
				KubeObjectValue:   value,
				KubeObjectKey:     key,
			}
			generateHelmTemplate(ingressObject, responseMap)
		}
	}

	// push the folder to GitHub
	filesToPush := []string{"values.yaml"}
	errWalk := filepath.Walk("helm/templates", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}

		relativePath := strings.Replace(path, "helm\\", "", 1)

		if !info.IsDir() {
			filesToPush = append(filesToPush, relativePath)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Could not file walk helm directory: %s", errWalk)
	}

	utils.PushToGitHub(gitWorkTree, gitRepo, filesToPush)
}

func generateHelmTemplate(kubeObj models.KubeComponent, responseMap map[string]interface{}) {
	// creates deployment.yaml file inside the helm/templates folder
	yamlFile, _ := os.Create(templatesFolderPath + "/" + kubeObj.KubeComponentType + ".yaml")

	// catches the deployment map inside the configs map
	var kubeObjectMap = kubeObj.KubeObjectValue.(map[string]interface{})

	// generates the path placeholder in the deployment helm chart
	generatedKubeObject := generatePlaceholders(kubeObjectMap, responseMap, "", kubeObj.KubeObjectKey)
	generatedKubeObject["apiVersion"] = kubeObj.APIVersion
	generatedKubeObject["kind"] = kubeObj.Kind

	// writes the deployment object to deployment.yaml file
	yamlData, _ := yaml.Marshal(&generatedKubeObject)
	_, err := yamlFile.Write(yamlData)
	if err != nil {
		log.Fatal(err)
	}
}

func generatePlaceholders(obj map[string]interface{}, responseObj map[string]interface{}, extraKey string, typeKey string) map[string]interface{} {
	for k, v := range obj {
		valueMap, isValueMap := v.(map[string]interface{}) // checks if v is of type map
		valueSlice, isValueSlice := v.([]interface{})      // checks if v is of type slice

		if isValueMap { // if v is a map recurse
			generatePlaceholders(valueMap, responseObj, extraKey+k+".", typeKey)
		} else if isValueSlice { // if v is a slice
			for _, s := range valueSlice { // iterates over the slice
				sliceMap, isSliceMap := s.(map[string]interface{}) // checks if s is a map
				if isSliceMap {                                    // if s is a map recurse
					generatePlaceholders(sliceMap, responseObj, extraKey+k+".", typeKey)
				}
			}
		} else { // if v is a leaf node set value for keys
			responseObj[k] = "{{.Values." + typeKey + "." + extraKey + k + "}}" // this temporarily holds the path
			obj[k] = "{{.Values." + typeKey + "." + extraKey + k + "}}"
		}
	}

	return obj
}
