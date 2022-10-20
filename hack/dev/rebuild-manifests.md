# Rebuilding Platform Manifests

## Rebuilding Image and CRDs

```bash
make docker-build docker-push extract IMG=robolaunchio/connection-hub-controller-manager:platform-v0.1.23
```

## Labeling Deployments

```bash
build=platform \
    make select-node \
        LABEL_KEY="robolaunch.io/organization" LABEL_VAL="robot-operator" && \
    make select-node \
        LABEL_KEY="robolaunch.io/department" LABEL_VAL="robotics" && \
    make select-node \
        LABEL_KEY="robolaunch.io/super-cluster" LABEL_VAL="aws-helsinki-1" && \
    make select-node \
        LABEL_KEY="robolaunch.io/cloud-instance" LABEL_VAL="instance1"
```