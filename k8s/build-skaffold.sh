#!/bin/bash
#BUILD OUR APP
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ./makako-gateway
#BUILD THE IMAGE
docker build -t registry.gitlab.com/jebo87/makako-gateway:test -f `pwd`/k8s/Dockerfile.dev .
#PUSH TO OUR REGISTRY
docker push registry.gitlab.com/jebo87/makako-gateway:test
#REMOVE UNUSED IMAGES
docker image prune -f