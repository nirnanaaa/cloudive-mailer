#!/bin/bash

apk add --update git
version=$(git describe --abbrev=1 --tags --always)
namespace="dev"

echo "Deploying version $version on namespace $namespace"

export TILLER_NAMESPACE=$namespace
export KUBECONFIG=$(pwd)/kubecfg
echo "$KUBECONF" |base64 -d |tee $KUBECONFIG
kubectl config set-context $(kubectl config current-context) --namespace=$namespace

helm upgrade thumb \
    --set "worker.image.tag=$version" \
    --set "master.image.tag=$version" \
    thumb-service