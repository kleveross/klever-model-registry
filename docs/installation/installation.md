# Installation

To use klever, you MUST install some dependent components as follows:
* [istio](https://github.com/istio/istio)
* [seldon-core](https://github.com/SeldonIO/seldon-core)
* [harbor](https://github.com/goharbor/harbor)

## Install istio
istio 的安装请参考 [istio 安装官方手册](https://istio.io/latest/docs/setup/install/)

## Install seldon core
According to Seldon Core official documentation, it support install by helm, please refer [install seldon-core by helm](https://github.com/SeldonIO/seldon-core/tree/master/helm-charts).

But for klever, we only support [istio](https://github.com/istio/istio), not support [ambassador](https://github.com/datawire/ambassador), and not support executor mode, since it MUST install by this manual, and then you can run klever-model-registry and klever-modejob-operator successfully.

### Install command
```bash
helm install seldon-core seldon-core-operator \
    --repo https://storage.googleapis.com/seldon-charts \
    --set usageMetrics.enabled=true \
    --namespace seldon-system \
    --set istio.enabled=true \
    --set istio.gateway=istio-system/kleveross-gateway \
    --set ambassador.enabled=false \
    --set executor.enabled=false \
    --set defaultUserID=0
```

## Install harbor
[Harbor](https://github.com/goharbor/harbor) is registry for the training model in klever-model-registry, please refer to the installation of Harbor [harbor-helm installation](https://github.com/goharbor/harbor-helm)

## Install klever
```bash
$ kubectl create namespace kleveross-system
$ git clone https://github.com/kleveross/klever-model-registry
$ cd klever-model-registry/manifests
$ helm install klever-model-registry ./klever-model-registry --namespace=kleveross-system --set ormb.domain={harbor address} --set model.externalAddress={model-registry-external-address}
$ helm install klever-modeljob-operator ./klever-modeljob-operator --namespace=kleveross-system --set ormb.domain={harbor address} --set model.externalAddress={model-registry-external-address}
```
