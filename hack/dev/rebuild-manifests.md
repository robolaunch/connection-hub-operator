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

## Deploying Operator

```bash
make apply
```

## Watching Controller Logs

```bash
connectionhub=logs k logs --follow pod/$(kubectl get pods -n connection-hub-system | tail -n1 | awk '{print $1}') -n connection-hub-system -c manager
```