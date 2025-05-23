sf-egress-ip-1184259
=============

This folder contains the code to [my answer][myanswer] to the "question"
[Use a static proxy IP for a pod in Kubernetes][question] on softwarerecs.stackexchange.com.

- [Prerequisites](#prerequisites)
- [Usage](#usage)
  - [Clone the repository](#clone-the-repository)
  - [Deploy without network policy](#deploy-without-network-policy)
  - [Deploy with network policy](#deploy-with-network-policy)

Prerequisites
-------------

* A git client
* [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl)
* [kustomize](https://kubectl.docs.kubernetes.io/installation/kustomize/) (optional, but helpful for adjustments)

Usage
-----

### Clone the repository

```
git clone https://github.com/mwmahlberg/stackoverflow-answers.git mwmahlberg-so-answers
cd mwmahlberg-so-answers/sf-egress-ip-1184259
```

### Deploy without network policy

```
kubectl apply -k base/
kubectl -n insights wait --for=jsonpath='{.status.loadBalancer.ingress}' svc/snmp-exporter
curl http://$(kubectl -n insights get svc snmp-exporter -o jsonpath='{.status.loadBalancer.ingress[0].ip}'):9116/metrics
```

### Deploy with network policy

```
kubectl apply -k overlays/with-policy
kubectl -n insights wait --for=jsonpath='{.status.loadBalancer.ingress}' svc/snmp-exporter
kubectl run -it --rm --image=alpine/curl --restart=Never debug \
-- ash -c "sleep 5 ; curl --connect-timeout 10 http://snmp-exporter.insights.svc.cluster.local:9116/metrics"
```

should give you the prometheus metrics whereas

```
curl -f http://$(kubectl -n insights get svc snmp-exporter -o jsonpath='{.status.loadBalancer.ingress[0].ip}'):9116/metrics
```

should result in `curl: (56) Recv failure: Connection reset by peer`.


[myanswer]: https://serverfault.com/a/1184272/238425
[question]: https://serverfault.com/questions/1184259/use-a-static-proxy-ip-for-a-pod-in-kubernetes