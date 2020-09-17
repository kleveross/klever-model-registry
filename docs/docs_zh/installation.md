# 安装

安装 klever 之前，需安装依赖的如下组件：
* [istio](https://github.com/istio/istio)
* [seldon-core](https://github.com/SeldonIO/seldon-core)
* [harbor](https://github.com/goharbor/harbor)

## 安装 istio
istio 的安装请参考 [istio 安装官方手册](https://istio.io/latest/docs/setup/install/)

## 安装 seldon core
Seldon core 官方支持 helm 安装，具体请参考 [install seldon-core by helm](https://github.com/SeldonIO/seldon-core/tree/master/helm-charts).

在 klever 中，当前对于流量分发只支持 [istio](https://github.com/istio/istio)， 暂不支持 [ambassador](https://github.com/datawire/ambassador), 并且暂不支持 seldon core 的 engine 模式，所以使用 klever 时安装 seldon core 必须设置一些额外的参数。

### 安装命令
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

## 安装 harbor
klever-model-registry 使用 [Harbor](https://github.com/goharbor/harbor) 存储训练模型，Harbor 的安装方式请参考 [harbor-helm 安装](https://github.com/goharbor/harbor-helm)

## 安装 klever
```bash
$ kubectl create namespace kleveross-system
$ git clone https://github.com/kleveross/klever-model-registry
$ cd klever-model-registry/manifests
$ helm install klever-model-registry ./klever-model-registry --namespace=kleveross-system --set ormb.domain={harbor address} --set model.externalAddress={model-registry-external-address}
$ helm install klever-modeljob-operator ./klever-modeljob-operator --namespace=kleveross-system --set ormb.domain={harbor address} --set model.externalAddress={model-registry-external-address}
```
