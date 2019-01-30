Prometheus Exporter Golang
---
Just a simple example how to export your custom metrics.

# Why this scaffolding

Basically this repo would like to show how to write a basic
producer/consumer application in _Go_ and create custom Prometheus metrics.
In fact most of the tasks are just faked, this means there is no an effective
processing, just some random data in order to simulate time consuming tasks
with failure as well.

# How to install on Minikube

```bash
$ kubectl create ns prometherion
namespace/prometherion created

$ eval $(minikube docker-env); docker build . -t prometherion/prometheus-exporter-golang:latest
Sending build context to Docker daemon  11.04MB
(...)
Successfully tagged prometherion/prometheus-exporter-golang:latest

$ kubectl -n prometherion create -f manifests/k8s
service/consumer-create-metrics created
service/consumer-delete-metrics created
service/consumer-read-metrics created
service/consumer-update-metrics created
deployment.extensions/consumer-create created
deployment.extensions/consumer-delete created
deployment.extensions/consumer-read created
deployment.extensions/consumer-update created
deployment.extensions/grafana created
service/grafana created
deployment.extensions/prometheus created
service/prometheus created
configmap/prometheus created
deployment.extensions/redis created
service/redis created
```

# View the dashboard (Grafana)

Basically I'm not using Ingress in order to approach _KISS_, so deploying Grafana
service as _NodePort_: you can grab it using this command.

```bash
$ k get svc grafana -o jsonpath='{.spec.ports[0].nodePort}'
3260
```

Obviously Minkube IP can be obtained with the following command:

```bash
$ minikube ip
192.168.99.100
```

Since it's a default installation, you can login with credentials
**admin** / **admin**.

So you can add the _Prometheus_ datasource using the address
`http://prometheus:9090` and load the
[dashboard](manifests/grafana/dashboard.json).

# How to produce some data

Just use the [produce.sh](produce.sh) script that create random actions and
push them as _Batch Jobs_.