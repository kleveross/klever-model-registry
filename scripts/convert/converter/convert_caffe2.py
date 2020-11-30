import os
import json
import argparse

import numpy as np
import onnx
import caffe2.python.onnx.frontend
from caffe2.proto import caffe2_pb2

from .base_convert import BaseConverter

INIT_NET = 'init_net.pb'
PREDICT_NET = 'predict_net.pb'
DEL_ATTR = 'ws_nbytes_limit'
MODEL_NAME = 'model.onnx'


def del_attr(netdef):
    for op in netdef.op:
        for i, attr in enumerate(op.arg):
            if attr.name == DEL_ATTR:
                op.arg.pop(i)


def np2onnx(s):
    def _modified2np(_s):
        if _s == 'float32':
            return 'float'
        if _s == 'float64':
            return 'double'
        return _s

    s = _modified2np(s)
    return onnx.TensorProto.DataType.Value(s.upper())


class Caffe2ToONNX(BaseConverter):
    def _load_model(self):
        self.init_net = self._find_with_name(INIT_NET)
        self.predict_net = self._find_with_name(PREDICT_NET)

    def _parse_input(self):
        value_info = {}
        for input in self.input_value:
            value_info[input['name']] = (np2onnx(input['dType']),
                                         tuple(input['size']))
        return value_info

    def _convert(self):
        self._load_model()
        value_info = self._parse_input()
        predict_net = caffe2_pb2.NetDef()
        with open(self.predict_net, 'rb') as f:
            predict_net.ParseFromString(f.read())

        init_net = caffe2_pb2.NetDef()
        with open(self.init_net, 'rb') as f:
            init_net.ParseFromString(f.read())

        del_attr(predict_net)

        out_path = os.path.join(self.output_dir, 'model', MODEL_NAME)

        onnx_model = caffe2.python.onnx.frontend.caffe2_net_to_onnx_model(
            predict_net,
            init_net,
            value_info,
        )

        onnx.save(onnx_model, out_path)


if __name__ == '__main__':
    convert = Caffe2ToONNX()
    convert.convert()
