klever-model-registry is an open-source model registry to manage machine learning model.

klever-model-registry helps you exract Machine Learning/Deep Learning models, or convert model from source format to anothner format.

## Components

* klever-model-registry: a API gateway for frontend.
* klever-model-operator: a CRD controller of ModelJob to manager extraction and conversion of model.

## Support model format for extract

* SavedModel
* ONNX
* GraphDef
* NetDef
* Keras H5
* CaffeModel
* TorchScript
* MXNetParams
* PMML 

## Support model format for convert

* MXNetParams to ONNX
* Keras H5 to SavedModel
* CaffeModel to NetDef 

## Compile from source

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

klever-model-registry project is part of Clever, a Cloud Native Machine Learning platform. We are going to open source a community edition soon.

The Clever slack workspace is caicloud-clever.slack.com. To join, click this [invitation to our Slack workspace](https://join.slack.com/t/caicloud-clever/shared_invite/zt-efz4rdrm-kcOg0Qvs_B8aIWGdZv9E6g).
