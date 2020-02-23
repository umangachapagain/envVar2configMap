package main

import (
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	const deploymentFile string = "./deployment.yaml"
	deployment := &appsv1.Deployment{}
	bytes, err := ioutil.ReadFile(deploymentFile)
	if err != nil {
		panic(err.Error())
	}
	err = yaml.Unmarshal(bytes, deployment)
	if err != nil {
		panic(err.Error())
	}

	resultFile := fmt.Sprintf("%s-config.yaml", deployment.Name)
	configMap := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
	}
	for _, container := range deployment.Spec.Template.Spec.Containers {
		configMap.Name = container.Name
		configMap.Namespace = deployment.Namespace
		configMap.Data = mapEnvVars(container.Env)
		result, err := yaml.Marshal(configMap)
		if err != nil {
			panic(err.Error())
		}
		err = ioutil.WriteFile(resultFile, result, 438)
		if err != nil {
			panic(err.Error())
		}
	}
}

func mapEnvVars(envs []corev1.EnvVar) map[string]string {
	result := make(map[string]string)
	if len(envs) < 1 {
		return result
	}
	for _, env := range envs {
		result[env.Name] = env.Value
	}
	return result
}
