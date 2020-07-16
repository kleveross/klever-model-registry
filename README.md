# Klever Model Registry

<a href="https://join.slack.com/t/kleveross/shared_invite/zt-g0eoiyq9-9OwiI7c__oV79bh_94MyTw">
    <img src="https://cdn.brandfolder.io/5H442O3W/as/pl546j-7le8zk-5guop3/Slack_RGB.png" alt="Slack" height =30px/></a>

English | [中文](./README_zh.md)

Klever Model Registry is a Cloud Native ML model registry. Use Klever Model Registry in order to:

- Manage your ML models
- Version and deliver your ML models with the existing infrastructures
- Keep track of ML models' hyperparameters and so on to help decision makers
- Convert models between different formats (e.g. TensorFlow SavedModel, ONNX)
- Serve the model (Coming Soon!)
- Get the standalone executable program to deploy ML model inference services on edge devices/servers (Coming Soon!)

Klever Model Registry's features:

- Deploy with Docker and Kubernetes
- Keep non-Invasive for your business
- Manage ML models like Docker (With the help of [kleveross/ormb](https://github.com/kleveross/ormb))
- Reuse Harbor to store models, without any new infrastructure
- Extract models signatures for:
    - SavedModel
    - ONNX
    - GraphDef
    - NetDef
    - Keras H5
    - CaffeModel
    - TorchScript
    - MXNetParams
    - PMML 
- Convert models from:
    - MXNetParams to ONNX
    - Keras H5 to SavedModel
    - CaffeModel to NetDef

## Installation

### Build the image

Clone:

```
$ git clone https://github.com/kleveross/klever-model-registry
$ cd klever-model-registry
```

Get the dependencies:

```
$ go mod tidy
```

Build:

```
$ make docker-build
```

### Install

To deploy klever-model-registry, execute following command:

```
$ kubectl apply -f https://raw.githubusercontent.com/kleveross/klever-model-registry/master/config/manager/klever-model-registry.yaml
$ kubectl apply -f https://raw.githubusercontent.com/kleveross/klever-model-registry/master/config/manager/klever-model-operator.yaml
```

## Community

klever-model-registry project is part of Klever, a Cloud Native Machine Learning platform.

The Klever slack workspace is klever.slack.com. To join, click this [invitation to our Slack workspace](https://join.slack.com/t/kleveross/shared_invite/zt-g0eoiyq9-9OwiI7c__oV79bh_94MyTw).
