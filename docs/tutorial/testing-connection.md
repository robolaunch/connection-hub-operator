# Testing Connection

In this document, it's aimed to test both Submariner and Federation connection.

## Testing Submariner

In cloud instance, create an `nginx` deployment. Exec into pod and ping physical instance.

```
kubectl create deployment nginx  --image nginx
kubectl exec -it pod/nginx-<POD-POSTFIX> -- bash
# inside nginx pod
apt-get update
apt-get install -y iputils-ping
ping 10.20.1.0 # get a valid IP from physical instance cluster CIDR
```

If ping process is successful, Submariner connection is opened between instances.

## Testing Federation

Since the cloud instance is federation host and the physical instance is federation member, it's enough to deploy a `FederatedNamespace` to cloud instance. (should federate namespace to both cloud instance and physical instance)

```
# in cloud instance
kubectl create ns connection-test
```

```yaml
# in cloud instance
# kubectl apply -f manifest.yaml
apiVersion: types.kubefed.io/v1beta1
kind: FederatedNamespace
metadata:
  name: connection-test
  namespace: connection-test
spec:
  placement:
    clusters:
    - name: <CLOUD-INSTANCE-NAME>
    - name: <PHYSICAL-INSTANCE-NAME>
```

Now check that if the namespace `connection-test` is created in the physical instance.

```
kubectl get ns connection-test
```

If the namespace is created, connection for federation is successful.