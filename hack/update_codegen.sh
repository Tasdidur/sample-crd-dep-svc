#!/bin/bash
vendor/k8s.io/code-generator/generate-groups.sh all \
	github.com/Tasdidur/xcrd/pkg/client \
	github.com/Tasdidur/xcrd/pkg/apis \
	xapi.com:v1
	#--go-header-file /home/office/go/src/github.com/Tasdidur/xcrd/vendor/k8s.io/code-generator/hack/boilerplate.go.txt

