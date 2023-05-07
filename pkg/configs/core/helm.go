package core

import (
	"config-generator/models"
	"config-generator/pkg/configs"
	"config-generator/pkg/utils"
	"encoding/base64"
	"fmt"
	"github.com/go-git/go-git/v5"
	"golang.org/x/exp/slices"
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

const (
	deployment  = "deploymentBased"
	statefulSet = "statefulSetBased"
	daemonSet   = "daemonSetBased"
	job         = "jobBased"
	cronJob     = "cronJobBased"
)

var templatesFolderPath = tmpHelmFolder + "/" + templatesFolder

// basic type is deployment based

func GenerateValueAndChartFiles(configs models.ConfigurationRequest, repoName string) (*git.Worktree, *git.Repository) {
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
	yamlDataValues, _ := yaml.Marshal(&configs)
	_, err = valuesFile.Write(yamlDataValues)
	if err != nil {
		log.Printf("Could not write data to values.yaml %s", err)
	}

	//creates chart.yaml file inside the helm folder
	chartFile, err := os.Create(tmpHelmFolder + "/Chart.yaml")
	if err != nil {
		log.Printf("Could not create Chart.yaml file: %s", err)
	}

	chart := models.Chart{
		APIVersion:  "v2",
		Name:        configs.AppName + "-chart",
		Version:     configs.Version,
		Description: "A Helm chart for " + configs.AppName + " application.",
	}

	yamlDataChart, _ := yaml.Marshal(chart)
	_, err = chartFile.Write(yamlDataChart)
	if err != nil {
		log.Printf("Could not write data to Chart.yaml %s", err)
	}

	return gitWorkTree, gitRepo
	//return nil, nil
}

func ConfigureTemplatesAndCharts(appDataRequest models.ConfigurationRequest, argoRepoName string, gitWorkTree *git.Worktree, gitRepo *git.Repository) (*git.Worktree, *git.Repository) {
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

	//microservices := appDataRequest.Microservices
	microservices := appDataRequest.Microservices
	for index, microservice := range microservices {

		microserviceName := microservice.ServiceName

		// {{- with (index .Values.grafana.ingress.hosts 0) }}
		// paths in the values.yaml to get placeholders
		//serviceNamePath := "microservices " + strconv.Itoa(index) + " .serviceName"
		//containerPortPath := "microservices " + strconv.Itoa(index) + " .containerPort"
		//avgReplicasPath := "microservices " + strconv.Itoa(index) + " .avgReplicas"
		//dockerImagePath := "microservices " + strconv.Itoa(index) + " .dockerImage"
		//minReplicasPath := "microservices " + strconv.Itoa(index) + " .minReplicas"
		//maxReplicasPath := "microservices " + strconv.Itoa(index) + " .maxReplicas"

		serviceNamePath := "microservices." + index + ".serviceName"
		//containerPortPath := "microservices." + index + ".containerPort"
		//avgReplicasPath := "microservices." + index + ".avgReplicas"
		dockerImagePath := "microservices." + index + ".dockerImage"
		//minReplicasPath := "microservices." + index + ".minReplicas"
		//maxReplicasPath := "microservices." + index + ".maxReplicas"
		//maxCPUPath := "microservices." + index + ".maxCPU"
		//maxMemoryPath := "microservices." + index + ".maxMemory"

		serviceNamePlaceholder := getPlaceholder(serviceNamePath)

		// configmaps
		configMap := models.ConfigMap{
			APIVersion: "v1",
			Kind:       "ConfigMap",
			Data:       make(map[string]string),
		}
		for iEnv, env := range microservices[index].Envs {
			configMap.Data[iEnv] = env.Value
		}
		configMap.Metadata.Name = serviceNamePlaceholder
		createManifestFile(configMap, microserviceName+"-configmap")

		// secret
		// fixme: duplicated envs in both secrets and configmaps
		// fixme: username and password should be saved in secrets not in configmaps
		secret := models.Secret{
			ApiVersion: "v1",
			Kind:       "Secret",
			Type:       "Opaque",
			Data:       make(map[string]string),
		}
		secret.Metadata.Name = serviceNamePlaceholder
		for iEnv, env := range microservices[index].Envs {
			// secret values should be encoded to base64
			secret.Data[iEnv] = string(encodeBase64(env.Value))
		}
		createManifestFile(secret, microserviceName+"-secret")

		// deploymentObj
		deploymentObj := models.Deployment{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		}

		deploymentObj.Metadata.Name = serviceNamePlaceholder
		deploymentObj.Metadata.Labels.App = serviceNamePlaceholder
		//deploymentObj.Spec.Replicas = getPlaceholder(avgReplicasPath)
		deploymentObj.Spec.Replicas = microservices[index].AvgReplicas
		deploymentObj.Spec.Selector.MatchLabels.App = serviceNamePlaceholder
		deploymentObj.Spec.Template.Metadata.Labels.App = serviceNamePlaceholder
		deploymentObj.Spec.Template.Metadata.Labels.App = serviceNamePlaceholder

		// envs array
		var env models.Env
		var envsSlice []models.Env
		var envNamePath string
		for iEnv, _ := range microservices[index].Envs {
			// fixme: MYSQL_ROOT_PASSWORD and MYSQL_DATABASE getting added to the env section of staefulset of the mysql and also the microservice deploymentObj
			//if !strings.EqualFold(envName.Name, "MYSQL_ROOT_PASSWORD") && !strings.EqualFold(envName.Name, "MYSQL_DATABASE") {
			//	envNamePath = "microservices[" + strconv.Itoa(index) + "].envs[" + strconv.Itoa(indexEnv) + "].name"
			//	env = models.Env{
			//		Name: getPlaceholder(envNamePath),
			//	}
			//	env.ValueFrom.ConfigMapKeyRef.Name = serviceNamePlaceholder
			//	env.ValueFrom.ConfigMapKeyRef.Key = getPlaceholder(envNamePath)
			//	envsSlice = append(envsSlice, env)
			//}
			envNamePath = "microservices." + index + ".envs." + iEnv + ".name" //todo: check env
			env = models.Env{
				Name: getPlaceholder(envNamePath),
			}
			env.ValueFrom.ConfigMapKeyRef.Name = serviceNamePlaceholder
			env.ValueFrom.ConfigMapKeyRef.Key = getPlaceholder(envNamePath)
			envsSlice = append(envsSlice, env)
		}

		// ports array
		//portDep := models.PortDep{ContainerPort: getPlaceholder(containerPortPath)}
		portDep := models.PortDep{ContainerPort: microservices[index].ContainerPort}
		container := models.Container{
			Name:  serviceNamePlaceholder,
			Image: getPlaceholder(dockerImagePath),
		}

		// adds ports array and envs array to the container struct of the deploymentObj struct
		container.Ports = append(container.Ports, portDep)
		container.Env = envsSlice

		// sets resource requests and limits for the container
		container.Resources.Requests.CPU = "100m"
		container.Resources.Requests.Memory = "500Mi"
		container.Resources.Limits.CPU = "1000m"
		container.Resources.Limits.Memory = "1Gi"

		deploymentObj.Spec.Template.Spec.Containers = append(deploymentObj.Spec.Template.Spec.Containers, container)
		createManifestFile(deploymentObj, microserviceName+"-deployment")

		// hpa
		hpa := models.HorizontalPodAutoscaler{
			APIVersion: "autoscaling/v1",
			Kind:       "HorizontalPodAutoscaler",
		}
		hpa.Metadata.Name = serviceNamePlaceholder
		hpa.Spec.ScaleTargetRef.APIVersion = "apps/v1"
		hpa.Spec.ScaleTargetRef.Kind = "Deployment"
		hpa.Spec.ScaleTargetRef.Name = serviceNamePlaceholder
		//hpa.Spec.MinReplicas = getPlaceholder(minReplicasPath)
		hpa.Spec.MinReplicas = microservices[index].MinReplicas
		//hpa.Spec.MaxReplicas = getPlaceholder(maxReplicasPath)
		hpa.Spec.MaxReplicas = microservices[index].MaxReplicas
		metrics := models.Metrics{
			Type: "Resource",
		}
		metrics.Resource.Name = "cpu"
		metrics.Resource.Target.Type = "Utilization"
		metrics.Resource.Target.AverageUtilization = 50
		hpa.Spec.Metrics = append(hpa.Spec.Metrics, metrics)
		createManifestFile(hpa, microserviceName+"-hpa")

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
		//createManifestFile(vpa, index+"-vpa") // fixme: argo error

		// service
		service := models.Service{
			ApiVersion: "v1",
			Kind:       "Service",
		}
		service.Metadata.Name = serviceNamePlaceholder
		service.Spec.Selector.App = serviceNamePlaceholder
		portSvc := models.PortSVC{
			Name: "http",
			Port: 80,
			//TargetPort: getPlaceholder(containerPortPath),
			TargetPort: microservices[index].ContainerPort,
		}
		service.Spec.Ports = append(service.Spec.Ports, portSvc)
		service.Spec.Type = "ClusterIP"
		createManifestFile(service, microserviceName+"-service")

		// ingress (rest)
		path := models.Path{
			Path:     "/" + serviceNamePlaceholder + "(/|$)(.*)",
			PathType: "Prefix",
		}
		path.Backend.Service.Name = serviceNamePlaceholder
		path.Backend.Service.Port.Name = "http"
		paths = append(paths, path)

		// If stateful service
		if slices.Contains(microservice.ServiceEvaluation.KubeConfigType, statefulSet) {
			if strings.EqualFold(appDataRequest.EvaluationResult.Database, "MySQL") {
				createMySQLConfigs(serviceNamePlaceholder, microserviceName)
			}
		}
		// If the microservice use MySQL
		//configTypes := microservice.Configs
		//isMysql := slices.Contains(configTypes, "mysql")
		//if isMysql {
		//	createMySQLConfigs(serviceNamePlaceholder, index)
		//}
	}

	rule := models.Rule{
		Host: getPlaceholder(hostNamePath),
		HTTP: struct {
			Paths []models.Path `yaml:"paths"`
		}{},
	}
	rule.HTTP.Paths = paths
	ingress.Spec.Rules = append(ingress.Spec.Rules, rule)
	createManifestFile(ingress, appDataRequest.AppName+"-ingress")

	// push the folder to GitHub
	pushHelmTemplatesToGitHub(gitWorkTree, gitRepo)

	// adds monitoring using kube-prometheus stack
	if appDataRequest.Monitoring == true {
		return configs.CreatePrometheusStackValuesYaml(appDataRequest.ClusterIPs, argoRepoName)
	}

	return nil, nil

}

func createMySQLConfigs(servicePlaceholder string, microserviceName string) {
	storage := "1Gi"
	localLabel := "local"
	numReplicas := 3
	statefulSetAppLabel := servicePlaceholder + "-mysql"

	// PV
	pv := models.PersistentVolume{
		ApiVersion: "v1",
		Kind:       "PersistentVolume",
		Spec: struct {
			StorageClassName string `yaml:"storageClassName"`
			Capacity         struct {
				Storage string `yaml:"storage"`
			} `yaml:"capacity"`
			AccessModes []string `yaml:"accessModes"`
			HostPath    struct {
				Path string `yaml:"path"`
			} `yaml:"hostPath"`
		}{},
	}

	pv.Metadata.Name = servicePlaceholder
	pv.Metadata.Labels.Type = localLabel
	pv.Spec.StorageClassName = "default"
	pv.Spec.Capacity.Storage = storage
	pv.Spec.AccessModes = []string{"ReadWriteOnce"}
	pv.Spec.HostPath.Path = "/mnt/data/mysql"
	createManifestFile(pv, microserviceName+"-pv")

	// PVC
	pvc := models.PersistentVolumeClaim{
		ApiVersion: "v1",
		Kind:       "PersistentVolumeClaim",
	}

	pvc.Metadata.Name = servicePlaceholder
	pvc.Metadata.Labels.App = servicePlaceholder
	pvc.Spec.AccessModes = []string{"ReadWriteOnce"}
	pvc.Spec.Resources.Requests.Storage = storage
	pvc.Spec.Selector.MatchLabels.Type = localLabel
	createManifestFile(pvc, microserviceName+"-pvc")

	statefulSetModel := models.StatefulSet{
		ApiVersion: "apps/v1",
		Kind:       "StatefulSet",
	}

	statefulSetModel.Metadata.Name = servicePlaceholder
	statefulSetModel.Spec.ServiceName = servicePlaceholder // a headless service will be created with this name
	statefulSetModel.Spec.Selector.MatchLabels.App = statefulSetAppLabel
	statefulSetModel.Spec.Replicas = numReplicas
	statefulSetModel.Spec.Template.Metadata.Labels.App = statefulSetAppLabel
	// ports
	port := models.PortDep{
		ContainerPort: 3306,
		Name:          "mysql",
	}
	ports := []models.PortDep{port}
	//envs
	envPW := models.Env{
		Name: "MYSQL_ROOT_PASSWORD",
	}
	envPW.ValueFrom.ConfigMapKeyRef.Name = servicePlaceholder
	envPW.ValueFrom.ConfigMapKeyRef.Key = "MYSQL_ROOT_PASSWORD"

	envDB := models.Env{
		Name: "MYSQL_DATABASE",
	}
	envDB.ValueFrom.ConfigMapKeyRef.Name = servicePlaceholder
	envDB.ValueFrom.ConfigMapKeyRef.Key = "MYSQL_DATABASE"

	envs := []models.Env{envPW, envDB}

	// volumeMounts
	volumeName := "mysql-data"
	volumePath := "/var/lib/mysql"
	volumeMount := models.VolumeMounts{
		Name:      volumeName,
		MountPath: volumePath,
	}
	// container
	container := models.Container{
		Name:         servicePlaceholder + "-mysql",
		Image:        "mysql:latest",
		Ports:        ports,
		Env:          envs,
		VolumeMounts: []models.VolumeMounts{volumeMount},
	}
	// volumeClaimTemplates
	volumeClaimTemplate := models.VolumeClaimTemplate{}
	volumeClaimTemplate.Metadata.Name = servicePlaceholder
	volumeClaimTemplate.Spec.AccessModes = []string{"ReadWriteOnce"}
	volumeClaimTemplate.Spec.Resources.Requests.Storage = storage
	volumeClaimTemplate.Spec.StorageClassName = "default"

	statefulSetModel.Spec.Template.Spec.Containers = []models.Container{container}
	statefulSetModel.Spec.VolumeClaimTemplates = []models.VolumeClaimTemplate{volumeClaimTemplate}
	volume := models.Volume{
		Name: volumeName,
		PersistentVolumeClaim: struct {
			ClaimName string `yaml:"claimName"`
		}{servicePlaceholder},
	}
	statefulSetModel.Spec.Template.Spec.Volumes = []models.Volume{volume}
	createManifestFile(statefulSetModel, microserviceName+"-statefulset")

	// service (to expose the stateful set)
	service := models.Service{
		ApiVersion: "v1",
		Kind:       "Service",
	}
	service.Metadata.Name = "mysql"
	service.Spec.Selector.App = statefulSetAppLabel
	portSvc := models.PortSVC{
		Protocol:   "TCP",
		Port:       3306,
		TargetPort: 3306,
	}
	service.Spec.Ports = []models.PortSVC{portSvc}
	createManifestFile(service, microserviceName+"-service-mysql")

}

func pushHelmTemplatesToGitHub(gitWorkTree *git.Worktree, gitRepo *git.Repository) {
	filesToPush := []string{"values.yaml", "Chart.yaml"} // add values.yaml anyway with charts inside helm/templates
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
