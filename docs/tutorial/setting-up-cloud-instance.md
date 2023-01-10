# Setting Up Cloud Instance

This document aims to give a tutorial about how robolaunch Connection Hub should be configured for **cloud instances**. 

## Installing Operator (Cloud Instance)
First of all, you need to install connection hub operator to your Kubernetes cluster. ([see prerequisites](https://github.com/robolaunch/connection-hub-operator/wiki/Tutorial-for-Physical-Instance-Cloud-Instance-Connection))

```
kubectl apply -f https://github.com/robolaunch/connection-hub-operator/releases/latest/download/connection-hub-operator.yaml
```

## Deploying Connection Hub (Cloud Instance)
After the operator is deployed, you need to deploy a `ConnectionHub` custom resource to make this cloud instance (Kubernetes cluster) the host of the connection hub. The inputs you need to set before deploying are listed below:

| Input               | Description                                                       |
|---------------------|-------------------------------------------------------------------|
| Cloud instance name | Name of the cloud instance. ([see here](https://github.com/robolaunch/connection-hub-operator/wiki/Configuration-Tips#naming-instances))                        |
| API Server URL      | API server URL of cluster.                                        |
| Cluster CIDR        | The CIDR pool used to assign IP addresses to pods in the cluster. ([see here](https://github.com/robolaunch/connection-hub-operator/wiki/Configuration-Tips#finding-cluster-cidr-of-instance)) |

Here is the manifest to setup connection hub. Beware that some of the inputs are inside angle brackets `<>` and they need to be populated:

```yaml
apiVersion: connection-hub.roboscale.io/v1alpha1
kind: ConnectionHub
metadata:
  labels:
    robolaunch.io/cloud-instance: <CLOUD-INSTANCE-NAME>
  name: connection-hub
spec:
  federationSpec:
    helmRepository:
        name: robolaunch
        url: "http://116.203.140.202:32401/robolaunch-helm-repository"
    helmChart:
      chartName: robolaunch/kubefed
      releaseName: kubefed
      version: 0.9.2
  submarinerSpec:
    helmRepository:
      name: robolaunch
      url: "http://116.203.140.202:32401/robolaunch-helm-repository"
    brokerHelmChart:
      chartName: robolaunch/submariner-k8s-broker
      releaseName: submariner-k8s-broker
      version: 0.6.0
    operatorHelmChart:
      chartName: robolaunch/submariner-operator
      releaseName: submariner-operator
      version: 0.10.1
    apiServerURL: <API-SERVER-URL>
    clusterCIDR: <CLUSTER-CIDR>
```

Watch connection hub resource phase until it's ready:

```
watch kubectl get connectionhub connection-hub
```

If the phase is `ReadyForOperation`, you have successfully configured the connection hub in your cloud instance. Next, you can connect your physical instance to cloud instance.