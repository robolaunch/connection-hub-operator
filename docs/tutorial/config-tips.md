This document refers to some data sources, conventions and rules while configuring a connection hub.

### Naming Instances

It's obligatory to use cluster domain names while naming instances. To learn cluster domain name, you can run:

```bash
kubectl get cm coredns -n kube-system -o jsonpath="{.data.Corefile}" \
  | grep ".local " \
  | awk -F ' ' '{print $2}'
```

Sample output will be:
```
cluster-one.local
```

It means that **YOU MUST** name that instance as `cluster-one`.

### Finding Cluster CIDR of Instance

It can be queried using:

```bash
kubectl get nodes <NODE-NAME> -o jsonpath='{.spec.podCIDR}'
```

### Finding Service CIDR of Instance

It can be queried using:

```bash
SVCRANGE=$(echo '{"apiVersion":"v1","kind":"Service","metadata":{"name":"tst"},"spec":{"clusterIP":"1.1.1.1","ports":[{"port":443}]}}' | kubectl apply -f - 2>&1 | sed 's/.*valid IPs is //')
echo $SVCRANGE
```