package kubectl

import (
	"fmt"
	"io/ioutil"
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

var KubeClients Config
var PreviousKubeconfig string

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
		log.Fatalf("config file not provided i.e %s env is unset", KUBECONFIG_MC)
	}

	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("failed to read the config file %s", err)
	}

	err = yaml.Unmarshal(yamlFile, &KubeClients)
	if err != nil {
		log.Fatalf("failed to unmarshal provided config %s", err)
	}

	PreviousKubeconfig = os.Getenv(KUBECONFIG)
}

func (c *Config) ExecCmd(args []string) {
	var err error
	for key, client := range c.Details {
		err = os.Setenv(KUBECONFIG, client.Kubeconfig)
		if err != nil {
			log.Fatalf("unable to switch context for %s", key)
		}
		fmt.Println("CLUSTER NAME: ", key)
		execKubectl(args)
		fmt.Println()
	}
	//defer rollBacKubeConfig()
}

func rollBacKubeConfig() {
	_ = os.Setenv("KUBECONFIG", PreviousKubeconfig)
}

func execKubectl(args []string) {
	out, err := exec.Command(KUBECTL, args...).CombinedOutput()
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Print(string(out))
}
