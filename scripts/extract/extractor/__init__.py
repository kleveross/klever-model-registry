import os
module = os.environ.get('EXTRACTOR','NULL')

if module=='onnx':
    from .extract_onnx import OnnxExtractor as Extractor
elif module=='caffe':
    from .extract_caffe import CaffeExtractor as Extractor
elif module=='netdef':
    from .extract_caffe2 import Caffe2Extractor as Extractor
elif module=='graphdef':
    from .extract_graphdef import TensorflowExtractor as Extractor
elif module=='h5':
    from .extract_keras import KerasExtractor as Extractor
elif module=='mxnetparams':
    from .extract_mxnet import MXNetExtractor as Extractor
elif module=='savedmodel':
    from .extract_savedmodel import TensorflowExtractor as Extractor
elif module=='torchscript':
    from .extract_torchscript import TorchscriptExtractor as Extractor
elif module=='pmml':
    from .extract_pmml import PMMLExtractor as Extractor
else:
    raise ImportError('module must be in one of [onnx, caffe, netdef, graphdef, h5, mxnetparams, savedmodel, torchscript]')