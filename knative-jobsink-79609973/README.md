knative-jobsink-79609973
========================

This folder contains the code to [my answer][myanswer] to the "question"
[Trigger knative jobSink from external source][question].

Usage
-----

> **Note**
>
> Differning from my answer, the ingress deployed heere
> is assumed to listen on 127.0.0.1:9090.
>
> If you have a different IP, you need to adjust [04_jobsink/ingress.yaml](./04_jobsink/ingress.yaml?plain=1#L8)

```plaintext
git clone https://github.com/mwmahlberg/stackoverflow-answers.git mwmahlberg-so-answers
cd mwmahlberg-so-answers/knative-jobsink-79609973/

kubectl apply -k 01_crds && \
  kubectl wait customresourcedefinitions.apiextensions \
  -l app.kubernetes.io/name=knative-eventing --for=condition=namesaccepted && \
  kubectl wait customresourcedefinitions.apiextensions \
  -l app.kubernetes.io/name=knative-serving --for=condition=namesaccepted

kubectl apply -k 02_serving && \
 kubectl -n knative-serving wait deployments \
 --for=condition=available --all --timeout=180s

kubectl apply -k 03_eventing && \
 kubectl -n knative-eventing wait deployments \
  --for=condition=available --all --timeout=180s

kubectl apply -k 04_jobsink && \
 kubectl -n demo wait --for=condition=ready jobsinks/job-sink-logger && \
 kubectl -n knative-eventing --for jsonpath='{.status.loadBalancer.ingress}' \
 ingress/jobsink-demo
```

With that setup, assuming your ingress connroller listens on 192.168.

[myanswer]: https://stackoverflow.com/a/79610285/1296707
[question]: https://stackoverflow.com/questions/79609973/trigger-knative-jobsink-from-external-source/79610285#79610285
