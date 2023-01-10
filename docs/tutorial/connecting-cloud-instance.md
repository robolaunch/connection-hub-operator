# Setting Up & Connecting Physical Instance

This document aims to give a tutorial about how robolaunch Connection Hub should be configured for **physical instances**. It's assumed that you have successfully configured connection hub in your cloud instance by following [this document](https://github.com/robolaunch/connection-hub-operator/wiki/1.-Setting-Up-Cloud-Instance).

## Registering Physical Instance (Cloud Instance)
To register physical instance to cloud instance, deploy the `PhysicalInstance` manifest below to your cloud instance. Notice that you need to set physical instance name before deploying ([see here](https://github.com/robolaunch/connection-hub-operator/wiki/Configuration-Tips#naming-instances)).
```yaml
# kubectl apply -f manifest.yaml
apiVersion: connection-hub.roboscale.io/v1alpha1
kind: PhysicalInstance
metadata:
  name: <PHYSICAL-INSTANCE-NAME>
```

## Querying Credentials (Cloud Instance)
To get credentials for connecting your physical instance, run the command below in your cloud instance:

```
# if you don't have yq installed, install it with `apt-get install -y yq`
# copy the YAML output
kubectl get connectionhub connection-hub -o jsonpath="{.status.connectionInterfaces.forPhysicalInstance}" | yq -P
```

## Installing Operator (Physical Instance)
In physical instance (k3s cluster), you need to install connection hub operator. ([see prerequisites for physical instance](https://github.com/robolaunch/connection-hub-operator/wiki/Tutorial-for-Physical-Instance-Cloud-Instance-Connection))

```
kubectl apply -f https://github.com/robolaunch/connection-hub-operator/releases/latest/download/connection-hub-operator.yaml
```

## Deploying Connection Hub (Physical Instance)
After the operator is deployed, you need to deploy a `ConnectionHub` custom resource to make this physical instance (k3s cluster) the member of the connection hub referencing a connection hub host (which is a cloud instance). The inputs you need to set before deploying are listed below:

|          Input         |                                       Description                                      |
|:----------------------:|:--------------------------------------------------------------------------------------:|
|   Cloud instance name  | Name of the cloud instance which is connection hub host. ([see here](https://github.com/robolaunch/connection-hub-operator/wiki/Configuration-Tips#naming-instances))                |
| Physical instance name | Name of the physical instance which will be connection hub member. ([see here](https://github.com/robolaunch/connection-hub-operator/wiki/Configuration-Tips#naming-instances))    |
|      Cluster CIDR      | The CIDR pool used to assign IP addresses to pods in the cluster. ([see here](https://github.com/robolaunch/connection-hub-operator/wiki/Configuration-Tips#finding-cluster-cidr-of-instance))     |
|      Service CIDR      | The CIDR pool used to assign IP addresses to services in the cluster. ([see here](https://github.com/robolaunch/connection-hub-operator/wiki/Configuration-Tips#finding-service-cidr-of-instance)) |

Here is the manifest to setup connection hub. `.spec` is completely gathered from the [Querying Credentials step](#querying-credentials-cloud-instance). Beware that some of the inputs are inside angle brackets `<>` and they need to be populated. (`<REDACTED>` fields contain secrets):

```yaml
# kubectl apply -f manifest.yaml
apiVersion: connection-hub.roboscale.io/v1alpha1
kind: ConnectionHub
metadata:
  labels:
    robolaunch.io/cloud-instance: <CLOUD-INSTANCE-NAME>
    robolaunch.io/physical-instance: <PHYSICAL-INSTANCE-NAME>
  name: connection-hub
spec:
  federationSpec:
    helmChart:
      chartName: robolaunch/kubefed
      releaseName: kubefed
      version: 0.9.2
    helmRepository:
      name: robolaunch
      url: http://116.203.140.202:32401/robolaunch-helm-repository
  instanceType: PhysicalInstance
  submarinerSpec:
    apiServerURL: <REDACTED>
    broker:
      ca: <REDACTED>
      token: <REDACTED>
    brokerHelmChart:
      chartName: robolaunch/submariner-k8s-broker
      releaseName: submariner-k8s-broker
      version: 0.6.0
    clusterCIDR: <PHYSICAL-INSTANCE-CLUSTER-CIDR>
    helmRepository:
      name: robolaunch
      url: http://116.203.140.202:32401/robolaunch-helm-repository
    operatorHelmChart:
      chartName: robolaunch/submariner-operator
      releaseName: submariner-operator
      version: 0.10.1
    presharedKey: <REDACTED>
    serviceCIDR: <PHYSICAL-INSTANCE-SERVICE-CIDR>
```

Watch connection hub resource phase until it's ready (`ReadyForOperation`):

```
watch kubectl get connectionhub connection-hub
```

## Setting Physical Instance Credentials (Cloud Instance)

The last step is to set physical instance credentials in cloud instance. First, check that if physical instance's phases match with the output below.

```
NAME      GATEWAY            HOSTNAME   CLUSTER ID   SUBNETS                           MULTICAST   FEDERATION              PHASE
robot01   ip-172-44-131-74   robot01    robot01      ["10.20.2.0/24","10.20.1.0/24"]   Connected   WaitingForCredentials   Registered
```

If the `Federation` field is `WaitingForCredentials`, update the manifest deployed in [Registering Physical Instance step](#registering-physical-instance-cloud-instance) as below:

|            Input           |                                    Description                                    |
|:--------------------------:|:---------------------------------------------------------------------------------:|
|   Physical instance name   | Name of the physical instance which will be connection hub member. ([see here](https://github.com/robolaunch/connection-hub-operator/wiki/Configuration-Tips#naming-instances)) |
|       API Server URL       | API server URL of physical instance.                              |
| Certificate Authority Data | ([see here]())                                                                    |
|   Client Certificate Data  | ([see here]())                                                                    |
|       Client Key Data      | ([see here]())                                                                    |

```yaml
# kubectl apply -f manifest.yaml
apiVersion: connection-hub.roboscale.io/v1alpha1
kind: PhysicalInstance
metadata:
  name: <PHYSICAL-INSTANCE-NAME>
spec:
  server: <PHYSICAL-INSTANCE-API-SERVER-URL>
  credentials:
    certificateAuthority: <PHYSICAL-INSTANCE-CERT-AUTHORITY>
    clientCertificate: <PHYSICAL-INSTANCE-CLIENT-CERT>
    clientKey: <PHYSICAL-INSTANCE-CLIENT-KEY>
```

After updating `PhysicalInstance` resource, you can watch resource phase until it's ready (`Connected`):

```
watch kubectl get physicalinstance <PHYSICAL-INSTANCE-NAME>
```

Example output:
```
NAME      GATEWAY            HOSTNAME   CLUSTER ID   SUBNETS                           MULTICAST   FEDERATION   PHASE
robot01   ip-172-44-131-74   robot01    robot01      ["10.20.2.0/24","10.20.1.0/24"]   Connected   Connected    Connected
```

If the both connections are successful, proceed the step [Testing Connection](https://github.com/robolaunch/connection-hub-operator/wiki/3.-Testing-Connection).