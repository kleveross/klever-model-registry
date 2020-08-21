# Klever Model Registry

[![Go Report Card](https://goreportcard.com/badge/github.com/kleveross/klever-model-registry)](https://goreportcard.com/report/github.com/kleveross/klever-model-registry)
[![Build Status](https://travis-ci.com/kleveross/klever-model-registry.svg?branch=master)](https://travis-ci.com/kleveross/klever-model-registry)
[![Coverage Status](https://coveralls.io/repos/github/kleveross/klever-model-registry/badge.svg?branch=master)](https://coveralls.io/github/kleveross/klever-model-registry?branch=master)

<a href="https://join.slack.com/t/kleveross/shared_invite/zt-g0eoiyq9-9OwiI7c__oV79bh_94MyTw">
    <img src="https://cdn.brandfolder.io/5H442O3W/as/pl546j-7le8zk-5guop3/Slack_RGB.png" alt="Slack" height =30px/></a>

[English](./README.md) | 中文


Klever Model Registry 是一个云原生的模型仓库。通过 Klever Model Registry，你可以：

- 管理你的机器学习模型的完整生命周期
- 利用已有的镜像仓库基础设施，对模型进行版本化和分发
- 跟踪模型的变化情况，查看模型的超参数等信息，帮助决策
- 对模型进行格式的转换，自动化地在 TensorFlow SavedModel，ONNX 等格式间转换，方便部署
- 利用模型，进行推理服务（Coming Soon!）
- 利用模型，获得适合边缘部署的独立二进制，便捷地进行边缘推理（Coming Soon!）

Klever Model Registry 的特性包括：

- 基于 Docker 与 Kubernetes 部署
- 对现有工作流的零侵入性
- 像 Docker 一样管理机器学习模型（在 [kleveross/ormb](https://github.com/kleveross/ormb) 的帮助下）
- 复用 Harbor 进行模型管理，不引入额外的基础设施
- 对以下格式的模型进行签名解析
    - SavedModel
    - ONNX
    - GraphDef
    - NetDef
    - Keras H5
    - CaffeModel
    - TorchScript
    - MXNetParams
    - PMML 
- 自动地进行模型格式间的转换（持续增加中）
    - MXNetParams 转为 ONNX
    - Keras H5 转为 SavedModel
    - CaffeModel 转为 NetDef 

## UI MockUp

<p align="center">
<img src="docs/images/model.png" height="300">
</p>

<p align="center">
<img src="docs/images/conversion.png" height="300">
</p>

## 安装

### 构建镜像

Clone:

```bash
$ git clone https://github.com/kleveross/klever-model-registry
$ cd klever-model-registry
```

Get the dependencies:

```bash
$ go mod tidy
```

Build:

```bash
$ make docker-build
```

### 部署

#### 安装 Harbor
klever-model-registry 使用 [Harbor](https://github.com/goharbor/harbor) 存储训练模型，Harbor 的安装方式请参考 [harbor-helm 安装](https://github.com/goharbor/harbor-helm)

#### 安装 klever-model-registry

```bash
$ kubectl create namespace kleveross-system
$ git clone https://github.com/kleveross/klever-model-registry
$ cd klever-model-registry/manifests
$ helm install klever-model-registry ./klever-model-registry --namespace=kleveross-system --set ormb.domain={harbor address}
$ helm install klever-modeljob-operator ./klever-modeljob-operator --namespace=kleveross-system --set ormb.domain={harbor address}
```

## 社区

klever-model-registry 是 Klever 云原生机器学习平台的子项目，Klever 的 Slack 是 klever.slack.com. 请利用[这一邀请链接](https://join.slack.com/t/kleveross/shared_invite/zt-g0eoiyq9-9OwiI7c__oV79bh_94MyTw)加入 Slack 讨论。
