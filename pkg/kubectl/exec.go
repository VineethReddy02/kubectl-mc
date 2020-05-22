package kubectl

import (
	"fmt"
	"io/ioutil"
	v1 "k8s.io/client-go/tools/clientcmd/api/v1"
	"log"
	"os"
	"os/exec"

	"gopkg.in/yaml.v2"
)

const (
	KUBECONFIG    = "KUBECONFIG"
	KUBECONFIG_MC = "KUBECONFIG_MC"
	KUBECTL       = "kubectl"
)

var ClientInfo ConfigDetails

type ConfigDetails struct {
	KubeMCClients Config
	KubeClients []ContextDetails
	KubeConfig bool
}

type ContextDetails struct {
	Name string
	Cluster string
}

type Config struct {
	Details map[string]Details `yaml:"details"`
}

type Details struct {
	Name       string `yaml:"name"`
	Kubeconfig string `yaml:"kubeconfig"`
}

func Initialise() {
	configFile := os.Getenv(KUBECONFIG_MC)
	if configFile == "" {
		ClientInfo.KubeConfig = true
		kubeConfig := os.Getenv(KUBECONFIG)
		if kubeConfig == "" {
			log.Fatal("failed to load clusters authentication details. provide clusters details either by setting env KUBECONFIG_MC or KUBECONFIG")
		}
		initialiseKubeConfig()
		return
	}

	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("failed to read the config file %s", err)
	}

	err = yaml.Unmarshal(yamlFile, &ClientInfo.KubeMCClients)
	if err != nil {
		log.Fatalf("failed to unmarshal provided config %s", err)
	}
}

// loads all the existing contexts from KUBECONFIG file.
func initialiseKubeConfig() {
	configFile := os.Getenv(KUBECONFIG)
	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("failed to read the config file %s", err)
	}

	config := &v1.Config{}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("failed to unmarshal provided config %s", err)
	}
	var context ContextDetails
	for _, d := range config.Contexts {
		context.Name = d.Name
		context.Cluster = d.Context.Cluster
		ClientInfo.KubeClients = append(ClientInfo.KubeClients, context)
	}
}

// sets the context between different clusters.
func (c *ConfigDetails) SetKubeContext(args []string) {
	var err error
	if !c.KubeConfig {
	for key, client := range c.KubeMCClients.Details {
		err = os.Setenv(KUBECONFIG, client.Kubeconfig)
		if err != nil {
			log.Fatalf("unable to switch context for %s", key)
		}
		fmt.Println("CLUSTER NAME: ", key)
		execKubectl(args)
		fmt.Println()
	}
	} else {
		for _, context := range c.KubeClients {
			_, err := exec.Command(KUBECTL, "config", "use-context", context.Name).CombinedOutput()
			if err != nil {
				log.Println(err.Error())
			}
			fmt.Println("CLUSTER NAME: ", context.Name)
			execKubectl(args)
			fmt.Println()
		}
	}
}

// executes the kubectl cmd.
func execKubectl(args []string) {
	out, err := exec.Command(KUBECTL, args...).CombinedOutput()
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Print(string(out))
}
