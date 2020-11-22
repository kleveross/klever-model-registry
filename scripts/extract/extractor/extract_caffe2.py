import os
import json
import collections

from caffe2.proto import caffe2_pb2
from .base_extract import BaseExtrctor
MODEL_TYPE = 'CAFFE2'
INIT_NET = 'init_net.pb'
PREDICT_NET = 'predict_net.pb'


class Caffe2Extractor(BaseExtrctor):
    def _extract_ops(self):
        op_types = map(lambda x: x.type, self.predict_net.op)
        ops = collections.Counter(op_types)
        return ops

    def _load_model(self):
        predict_path = self._find_with_name(PREDICT_NET)
        init_path = self._find_with_name(INIT_NET)
        self.predict_net = caffe2_pb2.NetDef()
        with open(predict_path, 'rb') as f:
            self.predict_net.ParseFromString(f.read())
