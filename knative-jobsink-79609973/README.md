knative-jobsink-79609973
========================

This folder contains the code to [my answer][myanswer] to the "question"
[Trigger knative jobSink from external source][question] on StackOverflow.

- [Usage](#usage)
  - [Setup](#setup)
  - [Submit a Job](#submit-a-job)


Usage
-----

> **Note**
>
> Differning from my answer, the ingress deployed heere
> is assumed to listen on 127.0.0.1:9090.
>
> If you have a different IP, you need to adjust [04_jobsink/ingress.yaml](./04_jobsink/ingress.yaml?plain=1#L8)

### Setup

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

### Submit a Job

With that setup, assuming your ingress connroller listens on 127.0.0.1:9090
(see above), you can trigger a Job like this:

```plaintext
$ curl -v \
  -H "content-type: application/json" \
  -H "ce-specversion: 1.0" \
  -H "ce-source: my/curl/command" \
  -H "ce-type: my.demo.event" \
  -H "ce-id: 123" \
  -d '{"details":"JobSinkDemo"}' \ 
  http://jobsink-demo.127-0-0-1.sslip.io:9090/demo/job-sink-logger

* Host jobsink-demo.127-0-0-1.sslip.io:9090 was resolved.
* IPv6: (none)
* IPv4: 127.0.0.1
*   Trying 127.0.0.1:9090...
* Connected to jobsink-demo.127-0-0-1.sslip.io (127.0.0.1) port 9090
* using HTTP/1.x
> POST /demo/job-sink-logger HTTP/1.1
> Host: jobsink-demo.127-0-0-1.sslip.io:9090
> User-Agent: curl/8.12.1
> Accept: */*
> content-type: application/json
> ce-specversion: 1.0
> ce-source: my/curl/command
> ce-type: my.demo.event
> ce-id: 123
> Content-Length: 25
> 
* upload completely sent off: 25 bytes
< HTTP/1.1 202 Accepted
< location: /namespaces/demo/name/job-sink-logger/sources/my/curl/command/ids/123
< date: Wed, 07 May 2025 11:18:19 GMT
< content-length: 0
< x-envoy-upstream-service-time: 27
< server: envoy
< 
* Connection #0 to host jobsink-demo.127-0-0-1.sslip.io left intact
```  

[myanswer]: https://stackoverflow.com/a/79610285/1296707
[question]: https://stackoverflow.com/questions/79609973/trigger-knative-jobsink-from-external-source/79610285#79610285
