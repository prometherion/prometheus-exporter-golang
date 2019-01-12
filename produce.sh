#!/usr/bin/env bash

declare -a ACTIONS=(
    "create"
    "read"
    "update"
    "delete"
)
INDEX=$( jot -r 1 0 $((${#ACTIONS[@]} - 1)) )
ACTION=${ACTIONS[INDEX]}

declare -a NAMES=(
    "foo"
    "bar"
    "bizz"
    "buzz"
)
INDEX=$( jot -r 1 0 $((${#NAMES[@]} - 1)) )
NAME=${NAMES[INDEX]}

declare -a IMAGES=(
    "redis:latest"
    "nginx:latest"
    "busybox:latest"
    "pause"
)
INDEX=$( jot -r 1 0 $((${#IMAGES[@]} - 1)) )
IMAGE=${IMAGES[INDEX]}

REPLICAS=`jot -r 1 1`

echo "apiVersion: batch/v1
kind: Job
metadata:
  name: ${ACTION}-${RANDOM}
  labels:
    app: prometheus-exporter-golang
spec:
  template:
    spec:
      containers:
      - name: app
        image: prometherion/prometheus-exporter-golang:latest
        imagePullPolicy: IfNotPresent
        command:
        - /prometheus-exporter-golang
        - --redis-connection-string=redis:6379
        - --redis-connection-tag=k8s
        - produce
        - app
        - --action=${ACTION}
        - --name=${NAME}
        - --image=${IMAGE}
        - --replicas=${REPLICAS}
      restartPolicy: Never
  backoffLimit: 1" | kubectl create -f -