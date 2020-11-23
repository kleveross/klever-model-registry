# Extract signature



## Usage

```python
EXTRACTOR={MODEL_TYPE} python extract.py -d {MODEL_PATH}
```

**MODEL_TYPE**  :  **[onnx, caffe, caffe2, graphdef, keras, mxnet, savedmodel, torchscript, tensorrt]**


## Framework Version
|Framework|version|
|-|-|
|Tensorflow|[TensorFlow 1.15.3](https://github.com/tensorflow/tensorflow/releases/tag/v1.15.3)|
|Pytorch|[1.7.0a0+6392713](https://github.com/pytorch/pytorch/commit/6392713)|
|MXNet|[1.6.0](https://github.com/apache/incubator-mxnet/releases/tag/1.6.0)|
|Caffe|[NvCaffe 0.17.3](https://github.com/NVIDIA/caffe/releases/tag/v0.17.3)|
|Caffe2|[1.7.0a0+6392713](https://github.com/pytorch/pytorch/commit/6392713)|
|Keras|[TensorFlow 1.15.3](https://github.com/tensorflow/tensorflow/releases/tag/v1.15.3)|
|ONNX|1.7.0|
|TensorRT|[TensorRT 7.1.3](https://docs.nvidia.com//deeplearning/tensorrt/release-notes/index.html)|
## Model format

```
Caffe
{MODEL_PATH}/
├── {NAME}.caffemodel
└── {NAME}.prototxt

caffe2
{MODEL_PATH}/
├── init_net.pb
└── predict_net.pb

mxnet
{MODEL_PATH}/
├── {NAME}.params
└── {NAME}-symbol.json

savedmodel
{MODEL_PATH}/
├── saved_model.pb
└── variables
    ├── variables.data-00000-of-00001
    └── variables.index

keras
{MODEL_PATH}/
   └──{NAME}.h5
  
onnx
{MODEL_PATH}/
└──{NAME}.onnx

torchscript
{MODEL_PATH}/
└── {NAME}.pt

graphdef
{MODEL_PATH}/
   └──{NAME}.graphdef

TensorRT
{MODEL_PATH}/
   └──{NAME}.[plan/engine]

```
