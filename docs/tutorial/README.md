# Tutorial for Physical Instance-Cloud Instance Connection

This tutorial explains connection hub setup. Here are the prerequisites for this tutorial:

- A cloud instance
  - Kubernetes cluster - version 1.19.X or above
    - `cert-manager` - version 1.8.0 or above
    - `coredns`
- A physical instance
  - k3s cluster - version 1.19.X or above
    - `cert-manager` - version 1.8.0 or above
    - `coredns`

## Labeling Nodes
For both instances, you should label the nodes that you will be working on.

### Label Cloud Instance Node
Node on the **cloud instance** should be labeled with cloud instance name ([refer this to name instance](https://github.com/robolaunch/connection-hub-operator/wiki/Configuration-Tips#naming-instances)) using key:
- `robolaunch.io/cloud-instance`

```
# in cloud instance
kubectl label node <NODE-NAME> robolaunch.io/cloud-instance=<CLOUD-INSTANCE-NAME>
```

### Label Physical Instance Node
Node on the **physical instance** should be labeled with cloud instance name and physical instance name ([refer this to name instance](https://github.com/robolaunch/connection-hub-operator/wiki/Configuration-Tips#naming-instances)) using keys:
- `robolaunch.io/cloud-instance`
- `robolaunch.io/physical-instance`

```
# in physical instance
kubectl label node <NODE-NAME> robolaunch.io/cloud-instance=<CLOUD-INSTANCE-NAME> robolaunch.io/physical-instance=<PHYSICAL-INSTANCE-NAME>
```

After this setup is done, you can proceed with [Setting Up Cloud Instance](https://github.com/robolaunch/connection-hub-operator/wiki/1.-Setting-Up-Cloud-Instance).