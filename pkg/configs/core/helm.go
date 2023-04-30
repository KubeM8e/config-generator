package core

import (
	"config-generator/models"
	"config-generator/pkg/utils"
	"encoding/base64"
	"fmt"
	"github.com/go-git/go-git/v5"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path/filepath"
	"strconv"
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

func GenerateValuesYamlFile(configs models.ConfigurationRequest, repoName string) (*git.Worktree, *git.Repository) {
	gitWorkTree, gitRepo := utils.CloneGitHubRepo(repoName, tmpHelmFolder)

	// create a helm folder
	//err := os.MkdirAll(tmpHelmFolder, os.ModePerm)
	//if err != nil {
	//	log.Printf("Could not make dir: %s", err)
	//}

	//creates values.yaml file inside the helm folder
	valuesFile, err := os.Create(tmpHelmFolder + "/values.yaml")
	if err != nil {
		log.Printf("Could not create values.yaml file: %s", err)
	}

	// converts map to yaml file
	yamlData, _ := yaml.Marshal(&configs)
	_, err = valuesFile.Write(yamlData)
	if err != nil {
		log.Printf("Could not write data to values.yaml %s", err)
	}

	return gitWorkTree, gitRepo
	//return nil, nil
}

func ConfigureHelmChart(configs models.ConfigurationRequest, gitWorkTree *git.Worktree, gitRepo *git.Repository) {
	appNamePath := "appName"
	hostNamePath := "hostName"

	// creates templates folder inside helm folder
	err := os.MkdirAll(templatesFolderPath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	// ingress
	ingress := models.Ingress{
		APIVersion: "networking.k8s.io/v1",
		Kind:       "Ingress",
		Metadata: struct {
			Name        string `yaml:"name"`
			Annotations struct {
				RewriteTarget string `yaml:"nginx.ingress.kubernetes.io/rewrite-target"`
			} `yaml:"annotations"`
		}{},
		Spec: struct {
			Rules []models.Rule `yaml:"rules"`
		}{},
	}
	ingress.Metadata.Name = getPlaceholder(appNamePath)
	ingress.Metadata.Annotations.RewriteTarget = "/$2"
	var paths []models.Path
	// rest of the ingress generation is inside the below loop - microservices details

	microservices := configs.Microservices
	for index, microservice := range microservices {

		// paths in the values.yaml to get placeholders
		serviceNamePath := "microservices[" + strconv.Itoa(index) + "].serviceName"
		containerPortPath := "microservices[" + strconv.Itoa(index) + "].containerPort"
		avgReplicasPath := "microservices[" + strconv.Itoa(index) + "].avgReplicas"
		dockerImagePath := "microservices[" + strconv.Itoa(index) + "].dockerImage"
		minReplicasPath := "microservices[" + strconv.Itoa(index) + "].minReplicas"
		maxReplicasPath := "microservices[" + strconv.Itoa(index) + "].maxReplicas"

		serviceNamePlaceholder := getPlaceholder(serviceNamePath)

		envs := microservice.Envs

		// configmaps
		configMap := models.ConfigMap{
			APIVersion: "v1",
			Kind:       "ConfigMap",
			Data:       make(map[string]string),
		}
		for _, env := range envs {
			configMap.Data[env.Name] = env.Value
		}
		configMap.Metadata.Name = serviceNamePlaceholder
		createManifestFile(configMap, microservice.ServiceName+"-configmap")

		// secret
		// fixme: duplicated envs
		secret := models.Secret{
			ApiVersion: "v1",
			Kind:       "Secret",
			Type:       "Opaque",
			Data:       make(map[string]string),
		}
		secret.Metadata.Name = serviceNamePlaceholder
		for _, env := range envs {
			// secret values should be encoded to base64
			secret.Data[env.Name] = string(encodeBase64(env.Value))
		}
		createManifestFile(secret, microservice.ServiceName+"-secret")

		// deployment
		deployment := models.Deployment{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
			Spec: struct {
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
						Containers []models.Container `yaml:"containers"`
					} `yaml:"spec"`
				} `yaml:"template"`
			}{},
		}

		deployment.Metadata.Name = serviceNamePlaceholder
		deployment.Metadata.Labels.App = serviceNamePlaceholder
		deployment.Spec.Replicas = getPlaceholder(avgReplicasPath)
		deployment.Spec.Selector.MatchLabels.App = serviceNamePlaceholder
		deployment.Spec.Template.Metadata.Labels.App = serviceNamePlaceholder
		deployment.Spec.Template.Metadata.Labels.App = serviceNamePlaceholder

		// envs array
		var env models.Env
		var envsSlice []models.Env
		var envNamePath string
		for indexEnv, _ := range envs {
			envNamePath = "microservices[" + strconv.Itoa(index) + "].envs[" + strconv.Itoa(indexEnv) + "].name"
			env = models.Env{
				Name: getPlaceholder(envNamePath),
			}
			env.ValueFrom.ConfigMapKeyRef.Name = serviceNamePlaceholder
			env.ValueFrom.ConfigMapKeyRef.Key = getPlaceholder(envNamePath)
			envsSlice = append(envsSlice, env)
		}

		// ports array
		portDep := models.PortDep{ContainerPort: getPlaceholder(containerPortPath)}
		container := models.Container{
			Name:  serviceNamePlaceholder,
			Image: getPlaceholder(dockerImagePath),
		}

		// adds ports array and envs array to the container struct of the deployment struct
		container.Ports = append(container.Ports, portDep)
		container.Env = envsSlice

		deployment.Spec.Template.Spec.Containers = append(deployment.Spec.Template.Spec.Containers, container)
		createManifestFile(deployment, microservice.ServiceName+"-deployment")

		// hpa
		hpa := models.HorizontalPodAutoscaler{
			APIVersion: "autoscaling/v1",
			Kind:       "HorizontalPodAutoscaler",
		}
		hpa.Metadata.Name = serviceNamePlaceholder
		hpa.Spec.ScaleTargetRef.APIVersion = "apps/v1"
		hpa.Spec.ScaleTargetRef.Kind = "Deployment"
		hpa.Spec.ScaleTargetRef.Name = serviceNamePlaceholder
		hpa.Spec.MinReplicas = getPlaceholder(minReplicasPath)
		hpa.Spec.MaxReplicas = getPlaceholder(maxReplicasPath)
		metrics := models.Metrics{
			Type: "Resource",
		}
		metrics.Resource.Name = "cpu"
		metrics.Resource.Target.Type = "Utilization"
		metrics.Resource.Target.AverageUtilization = 50
		hpa.Spec.Metrics = append(hpa.Spec.Metrics, metrics)
		createManifestFile(hpa, microservice.ServiceName+"-hpa")

		// VPA
		vpa := models.VerticalPodAutoscaler{
			APIVersion: "autoscaling.k8s.io/v1",
			Kind:       "VerticalPodAutoscaler",
		}
		vpa.Metadata.Name = serviceNamePlaceholder
		vpa.Spec.TargetRef.APIVersion = "apps/v1"
		vpa.Spec.TargetRef.Kind = "Deployment"
		vpa.Spec.TargetRef.Name = serviceNamePlaceholder
		vpa.Spec.UpdatePolicy.UpdateMode = "Auto"
		createManifestFile(vpa, microservice.ServiceName+"-vpa")

		// service
		service := models.Service{
			ApiVersion: "v1",
			Kind:       "Service",
		}
		service.Metadata.Name = serviceNamePlaceholder
		service.Spec.Selector.App = serviceNamePlaceholder
		portSvc := models.PortSVC{
			Name:       "http",
			Port:       80,
			TargetPort: getPlaceholder(containerPortPath),
		}
		service.Spec.Ports = append(service.Spec.Ports, portSvc)
		service.Spec.Type = "ClusterIP"
		createManifestFile(service, microservice.ServiceName+"-service")

		// ingress (rest)
		path := models.Path{
			Path:     "/" + serviceNamePlaceholder + "(/|$)(.*)",
			PathType: "Prefix",
		}
		path.Backend.Service.Name = serviceNamePlaceholder
		path.Backend.Service.Port.Name = "http"
		paths = append(paths, path)

	}

	rule := models.Rule{
		Host: getPlaceholder(hostNamePath),
		HTTP: struct {
			Paths []models.Path `yaml:"paths"`
		}{},
	}
	rule.HTTP.Paths = paths
	ingress.Spec.Rules = append(ingress.Spec.Rules, rule)
	createManifestFile(ingress, configs.AppName+"-ingress")

	// push the folder to GitHub
	pushHelmTemplatesToGitHub(gitWorkTree, gitRepo)

}

func pushHelmTemplatesToGitHub(gitWorkTree *git.Worktree, gitRepo *git.Repository) {
	filesToPush := []string{"values.yaml"} // add values.yaml anyway with charts inside helm/templates
	errWalk := filepath.Walk("helm/templates", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("File walk error: %s", err)
			return err
		}

		// the path is like helm/test/testFile - should remove helm part
		relativePath := strings.Replace(path, "helm\\", "", 1)

		if !info.IsDir() {
			filesToPush = append(filesToPush, relativePath)
		}

		return nil
	})

	if errWalk != nil {
		fmt.Printf("Could not file walk helm directory: %s", errWalk)
	}

	utils.PushToGitHub(gitWorkTree, gitRepo, filesToPush)
}

func encodeBase64(str string) []byte {
	dst := make([]byte, base64.StdEncoding.EncodedLen(len([]byte(str))))
	base64.StdEncoding.Encode(dst, []byte(str))
	return dst
}

func createManifestFile(manifest interface{}, fileName string) {
	yamlFile, errCreate := os.Create(templatesFolderPath + "/" + fileName + ".yaml")
	if errCreate != nil {
		fmt.Printf("Could not create yaml: %s", errCreate)
	}

	yamlData, _ := yaml.Marshal(manifest)
	_, errWrite := yamlFile.Write(yamlData)
	if errWrite != nil {
		log.Printf("Could not write the manifest file: %s", errWrite)
	}
}

func getPlaceholder(path string) string {
	segments := strings.Split(path, ".")
	segments[0] = "{{.Values." + segments[0]
	lastIndex := len(segments) - 1
	segments[lastIndex] = segments[lastIndex] + "}}"
	placeholder := strings.Join(segments, ".")
	return placeholder
}

// Old version -  refer to API /helm

func GenerateValuesYamlFile2(configs map[string]interface{}, repoName string) (*git.Worktree, *git.Repository) {
	gitWorkTree, gitRepo := utils.CloneGitHubRepo(repoName, tmpHelmFolder)

	// create a helm folder
	//err := os.MkdirAll(tmpHelmFolder, os.ModePerm)
	//if err != nil {
	//	log.Printf("Could not make dir: %s", err)
	//}

	//creates values.yaml file inside the helm folder
	valuesFile, err := os.Create(tmpHelmFolder + "/values.yaml")
	if err != nil {
		log.Printf("Could not create values.yaml file: %s", err)
	}

	// converts map to yaml file
	yamlData, _ := yaml.Marshal(&configs)
	_, err = valuesFile.Write(yamlData)
	if err != nil {
		log.Printf("Could not write data to values.yaml %s", err)
	}

	return gitWorkTree, gitRepo
}

func ConfigureHelmChart2(configs map[string]interface{}, gitWorkTree *git.Worktree, gitRepo *git.Repository) {

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
	// add values.yaml anyway with charts inside helm/templates
	filesToPush := []string{"values.yaml"}
	errWalk := filepath.Walk("helm/templates", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("File walk error: %s", err)
			return err
		}

		// the path is like helm/test/testFile - should remove helm part
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
