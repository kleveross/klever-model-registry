#!/bin/bash

############# Version Information Begin #############
# kubebernets 1.19 is test ok.
# istio version: v1.8.0
# seldon core: v1.4.0
# harbor: v2.1.1
# klever 0.1.0
############# Version Information End ###############

# Set it as k8s master ip.
export MASTER_IP=$1

# Set harbor NodePort port.
export HARBOR_PORT=30022

# Set klever-model-registry NodePort port.
export KLEVER_MODEL_REGISTRY_PORT=30100
export KLEVER_WEB_PORT=30200

# 
# Go to manifests directory, it is workdir.
CWD=$(pwd)

#
# Install istio, please reference https://istio.io/latest/docs/setup/install/helm/

curl -L https://istio.io/downloadIstio | ISTIO_VERSION=1.8.0 sh -
export PATH=$PWD/bin:$PATH
kubectl create namespace istio-system
# If you enabled third party tokens, the jwtPolicy option can be deleted
helm install --namespace istio-system istio-base istio-1.8.0/manifests/charts/base \
    --set global.jwtPolicy=first-party-jwt
helm install --namespace istio-system istiod istio-1.8.0/manifests/charts/istio-control/istio-discovery \
    --set global.hub="docker.io/istio" --set global.tag="1.8.0" \
    --set global.jwtPolicy=first-party-jwt
helm install --namespace istio-system istio-ingress istio-1.8.0/manifests/charts/gateways/istio-ingress \
    --set global.hub="docker.io/istio" --set global.tag="1.8.0" \
    --set global.jwtPolicy=first-party-jwt

#
# Install harbor
#
git clone https://github.com/kleveross/klever-model-registry.git
kubectl create namespace harbor-system

helm install harbor $CWD/klever-model-registry/manifests/harbor \
    --set expose.nodePort.ports.http.nodePort=$HARBOR_PORT \
    --set expose.type=nodePort \
    --set expose.tls.enabled=false \
    --set trivy.enabled=false \
    --set notary.enabled=false \
    --set persistence.enabled=false \
    --set trivy.ignoreUnfixed=true \
    --set trivy.insecure=true \
    --set externalURL=http://$MASTER_IP:$HARBOR_PORT \
    --set harborAdminPassword="ORMBtest12345" \
    --namespace harbor-system

#
# Install seldon, please reference https://docs.seldon.io/projects/seldon-core/en/latest/charts/seldon-core-operator.html
#
git clone -b klever-v1.5.0-release --single-branch https://github.com/kleveross/seldon-core.git
kubectl create namespace seldon-system
helm install seldon-core $CWD/seldon-core/helm-charts/seldon-core-operator \
    --set usageMetrics.enabled=true \
    --set istio.enabled=true \
    --set istio.gateway=istio-system/kleveross-gateway \
    --set ambassador.enabled=false \
    --set executor.enabled=false \
    --set defaultUserID=0 \
    --set image.registry=ghcr.io \
    --set image.repository=kleveross/seldon-core-operator \
    --set image.tag=0.1.0 \
    --namespace seldon-system

#
# Install Klever
#
kubectl create namespace kleveross-system

# Install Klever-modeljob-operator
helm install klever-modeljob-operator $CWD/klever-model-registry/manifests/modeljob-operator \
    --namespace=kleveross-system

# Install Klever-model-registry
helm install klever-model-registry $CWD/klever-model-registry/manifests/model-registry \
    --set externalAddress=$MASTER_IP:$KLEVER_MODEL_REGISTRY_PORT \
    --set service.nodePort=$KLEVER_MODEL_REGISTRY_PORT \
    --namespace=kleveross-system

# Install Klever-web
git clone https://github.com/kleveross/klever-web
helm install klever-web $CWD/klever-web/manifests/klever-web \
    --namespace=kleveross-system \
    --set service.nodePort=$KLEVER_WEB_PORT
    --set model.registry.address=http://$MASTER_IP:$KLEVER_MODEL_REGISTRY_PORT
