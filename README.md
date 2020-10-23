# Klever Model Registry

[![Go Report Card](https://goreportcard.com/badge/github.com/kleveross/klever-model-registry)](https://goreportcard.com/report/github.com/kleveross/klever-model-registry)
[![Build Status](https://github.com/kleveross/klever-model-registry/workflows/UnitTest/badge.svg
)](https://github.com/kleveross/klever-model-registry/actions?query=workflow%3AUnitTest)
[![Build Status](https://github.com/kleveross/klever-model-registry/workflows/E2ETest/badge.svg
)](https://github.com/kleveross/klever-model-registry/actions?query=workflow%3AE2ETest)
[![Coverage Status](https://coveralls.io/repos/github/kleveross/klever-model-registry/badge.svg?branch=master)](https://coveralls.io/github/kleveross/klever-model-registry?branch=master)

<a href="https://join.slack.com/t/kleveross/shared_invite/zt-g0eoiyq9-9OwiI7c__oV79bh_94MyTw">
    <img src="https://cdn.brandfolder.io/5H442O3W/as/pl546j-7le8zk-5guop3/Slack_RGB.png" alt="Slack" height =30px/></a>

English | [中文](./README_zh.md)

Klever Model Registry is a Cloud Native ML model registry. Use Klever Model Registry in order to:

- Manage your ML models
- Version and deliver your ML models with the existing infrastructures
- Keep track of ML models' hyperparameters and so on to help decision makers
- Convert models between different formats ( e.g. TensorFlow SavedModel, ONNX )
- Serve the model
- Get the standalone executable program to deploy ML model inference services on edge devices/servers ( Coming Soon! )

Klever Model Registry's features:

- Deploy with Docker and Kubernetes
- Keep non-Invasive for your business
- Manage ML models like Docker ( With the help of [kleveross/ormb](https://github.com/kleveross/ormb) )
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

See our [official documentations](/docs/README.md) for more information。

## UI MockUp

<p align="center">
<img src="docs/images/model.png" height="400">
</p>

<p align="center">
<img src="docs/images/conversion.png" height="400">
</p>

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

Please have a look at [docs/installation.md](docs/installation.md).

If you want to trial quickly, you can run installation script as follow.

```bash
$ wget https://raw.githubusercontent.com/kleveross/klever-model-registry/master/scripts/installation/install.sh
$ bash install.sh <master-ip>
```

## Community

klever-model-registry project is part of Klever, a Cloud Native Machine Learning platform.

The Klever slack workspace is klever.slack.com. To join, click this [invitation to our Slack workspace](https://join.slack.com/t/kleveross/shared_invite/zt-g0eoiyq9-9OwiI7c__oV79bh_94MyTw).
