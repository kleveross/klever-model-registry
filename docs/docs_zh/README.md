# Klever Model Registry

[English](../README.md) | 中文

Klever Model Registry 是一个使用 [Harbor](https://github.com/goharbor/harbor) 存储训练模型的云原生模型仓库。用户可以使用 Klever Model Registry 进行模型管理、模型解析、模型转换、模型服务。
Klever Model Registry 是一个开源项目，遵循 [Apache 2.0 开源协议](https://www.apache.org/licenses/LICENSE-2.0)，属于 Klever 云原生机器学习平台的子项目。以下简称 Klever。

<!-- // TODO:
## 入门指南

想要了解 Klever Model Registry，可以先完成我们的 [入门教程]() 并且查看我们的 [使用示例]()。
-->

## 深入理解

阅读以下文档以及我们的 [API 文档](https://kleveross.github.io/klever-model-registry/api/) 来学习如何使用 Klever Model Registry:

- [管理模型](##管理模型)
- [模型解析](##模型解析)
- [模型转换](##模型转换)
- [模型服务](##模型服务)

## 管理模型

用户可以通过调用 API 的方式上传和下载其模型，Klever 会使用 `ormb` 推送模型到 Harbor 中。

> `ormb` 是 Klever 云原生机器学习平台的另一开源项目，是一个 OCI(Open Container Initiative)-Based 的仓库，致力于利用已有的镜像仓库进行分发模型和模型的元数据。详见项目 [ORMB](https://github.com/kleveross/ormb)。

通过指定模型在 Harbor 中所存储的项目名，模型的名字及版本，即可将模型包上传至 Harbor 中。模型包需满足 `ormb` 规范，必须要有 `ormbfile.yaml`，在这个 yaml 文件中会存放模型的一些信息，如框架、格式等（后续我们将支持自动生成 `ormbfile.yaml`，Coming Soon!）。Klever 代理了 Harbor 的所有请求，若暂无 Harbor 项目，可以通过 Klever 创建项目。

## 模型解析

当模型被推送至 Harbor 完成后，Klever 会自动创建 `ModelJob` 进行模型解析。`ModelJob` 是 Klever 定义的一个 [CRD](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)， Klever 会根据被上传的模型格式填写 `ModelJob.Spec.Extraction`，并生成 Job 执行模型解析。目前支持的模型解析格式有：
  - SavedModel
  - ONNX
  - GraphDef
  - NetDef
  - Keras H5
  - CaffeModel
  - TorchScript
  - MXNetParams
  - PMML 

模型的具体解析过程由 `ModelJob` 生成并控制的 `Job` 的镜像来完成，在解析完毕后会生成更新后的 `ormbfile.yaml` 并推送到 `Harbor`。镜像中的解析脚本代码详见 [extract](/scripts/extract/extract.py)。

## 模型转换

目前 Klever 支持三种类型的模型格式转换：
  - MXNetParams 转为 ONNX
  - Keras H5 转为 SavedModel
  - CaffeModel 转为 NetDef

同样借助于 `ModelJob` 这个 `CRD`，用户可以通过调用 API 的形式来创建用于模型转换的 `ModelJob`。通过指定 `ModelJob.Spec.Conversion.Mmdnn.From` 和 `ModelJob.Spec.Conversion.Mmdnn.To` 来确定模型转换的原格式和目标格式。模型的具体转换过程由 `ModelJob` 生成并控制的 `Job` 的镜像来完成，在转换完毕后会生成更新后的 `ormbfile.yaml` 并推送到 `Harbor`。镜像中的转换脚本代码详见 [convert](/scripts/convert/base_convert/base_convert.py)。

## 模型服务

Klever 基于 [Seldon-Core](https://github.com/SeldonIO/seldon-core) 实现模型服务，创建模型服务会首先创建一个 `Seldon Deployment`，并在其 `Init Container` 中通过 [ormb-storage-initializer](https://github.com/kleveross/ormb/blob/master/build/ormb-storage-initializer/Dockerfile) 下载模型。若模型为 `PMML` 格式，将使用 [OpenScoring 镜像](/build/serving/openscoring/Dockerfile) 启动服务；若模型为其他 [Triton Server](https://docs.nvidia.com/deeplearning/triton-inference-server/master-user-guide/docs/model_repository.html#framework-model-definition) 支持的模型格式，将使用 [Triton Server 镜像](/build/serving/tensorrt/Dockerfile) 启动服务，镜像中会自动通过 `ormbfile.yaml` 中的信息生成 `Triton Server` 所需要的 [config.pbtxt](https://docs.nvidia.com/deeplearning/triton-inference-server/user-guide/docs/model_configuration.html#) 文件。

## 成为贡献者之一

如果您愿意为 Klever Model Registry 项目做出贡献，请参阅我们的 [贡献者指南](/CONTRIBUTING.md)。

