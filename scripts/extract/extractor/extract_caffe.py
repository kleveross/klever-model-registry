import os
os.environ['GLOG_minloglevel'] = '2'
import json
import collections

import caffe

from .base_extract import BaseExtrctor

MODEL_TYPE = 'CAFFE'
MODEL_EXTENSION = '.prototxt'
WEIGHTS_EXTENSION = '.caffemodel'


class CaffeExtractor(BaseExtrctor):
    def _extract_inputs(self):
        inputs = []
        for input_name in self.model.inputs:
            blob_data = self.model.blobs[input_name].data
            origin_inputs = {}
            origin_inputs['name'] = input_name
            origin_inputs['dType'] = blob_data.dtype.name
            origin_inputs['size'] = list(blob_data.shape)
            inputs.append(origin_inputs)
        return inputs

    def _extract_outputs(self):
        outputs = []
        for output_name in self.model.outputs:
            blob_data = self.model.blobs[output_name].data
            origin_outputs = {}
            origin_outputs['name'] = output_name
            origin_outputs['dType'] = blob_data.dtype.name
            origin_outputs['size'] = list(blob_data.shape)
            outputs.append(origin_outputs)
        return outputs

    def _extract_ops(self):
        op_types = map(lambda x: x.type, self.model.layers)
        ops = collections.Counter(op_types)
        return ops

    def _load_model(self):
        model_path = self._find_with_extension(MODEL_EXTENSION)
        weight_path = self._find_with_extension(WEIGHTS_EXTENSION)
        self.model = caffe.Net(model_path, caffe.TEST)
