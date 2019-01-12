build:
	eval $$(minikube docker-env); docker build . -t prometherion/prometheus-exporter-golang:latest
