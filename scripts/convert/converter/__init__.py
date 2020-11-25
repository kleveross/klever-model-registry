import os
srcFormat = os.environ.get('SOURCE_FORMAT', 'NULL')
dstFormat = os.environ.get('FORMAT', 'NULL')
if dstFormat == 'ONNX':
    if srcFormat == 'MXNetParams':
        from .convert_mxnet import MXNetToONNX as Converter
    elif srcFormat == 'NetDef':
        from .convert_caffe2 import Caffe2ToONNX as Converter
elif dstFormat == 'NetDef':
    from .convert_caffe import CaffeToCaffe2 as Converter
elif dstFormat == 'SavedModel':
    from .convert_keras import KerasToTensorFlow as Converter
else:
    raise ImportError(
        'SOURCE_FORMAT must in [ MXNetParams, NetDef, CaffeModel, H5 ],FORMAT must in [ ONNX, NetDef, SavedModel ]'
    )
