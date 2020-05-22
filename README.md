# kubectl-mc

kubectl plugin to support multiple kubernetes clusters in single go.

The motive behind this plugin is to manage set of kubernetes clusters with a kubectl plugin. All the kubectl actions can be performed using this plugin. The plugin cascades the performed operation to all the k8s clusters.

We access multiple clusters using the ```KUBECONFIG``` file by reading all the contexts. For this to work all you need is ```KUBECONFIG``` env set.
More details on configuring access to multiple clusters using a single ```KUBECONFIG``` file can be found [here](https://kubernetes.io/docs/tasks/access-application-cluster/configure-access-multiple-clusters/)

OR

You can provide the paths to the multiple ```KUBECONFIG``` files for this to work. As specified in the below configuration file.

Configuration:
```
details:
  # Desired cluster name
  K8S:
    # Path th the KUBECONFIG file
    kubeconfig: "/home/vineeth/Downloads/raw51.yaml"
  AKS:
    kubeconfig: "/home/vineeth/Downloads/aks52.yaml"
```
Set KUBECONFIG_MC env with above yaml file i.e export ```KUBECONFIG_MC=config.yaml```

**Note:** When both ```KUBECONFIG``` & ```KUBECONFIG_MC``` env are set ```KUBECONFIG_MC``` the custom env designed for this plugin to ease adding details on multiple ```KUBECONFIG``` file paths will take the precedence.
 
### Implementation details:

So, I didn't want to re-invent the wheel by supporting all commands & flags of ```kubectl``` cli. So to make implementation as simple as possible and to support all the features of kubectl. Plugin performs the kubectl operations by switching the context between multiple k8s clusters.  

### Usage:

Creating a pod in multiple clusters with single command.
```
 $ kubectl mc run --generator=run-pod/v1 vi --image vineeth0297/languages
CLUSTER NAME:  K8S
pod/vi created

CLUSTER NAME:  AKS
pod/vi created
 
```

Listing the pods from multiple clusters.

```
$ kubectl mc get pods
CLUSTER NAME:  K8S
NAME   READY   STATUS    RESTARTS   AGE
vi     1/1     Running   0          117s

CLUSTER NAME:  AKS
NAME            READY   STATUS    RESTARTS   AGE
aks-ssh-qxs2c   1/1     Running   0          8m39s
vi              1/1     Running   0          117s
```

This plugin is 100% compatible with kubectl commands & flags.

Feel free to share the feedback/thoughts/suggestions to improve this plugin in managing multiple clusters through cli.