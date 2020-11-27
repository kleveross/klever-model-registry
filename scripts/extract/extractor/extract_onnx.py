import os
import json
import collections

import onnx
import onnxruntime as onnxr
from .base_extract import BaseExtrctor

MODEL_TYPE = 'ONNX'
EXTENSION = '.onnx'


def modified2np(s):
    if s == 'float':
        return 'float32'
    if s == 'double':
        return 'float64'
    return s


modefied_dynamic = lambda x: -1 if isinstance(x, str) else x


class OnnxExtractor(BaseExtrctor):
    def _extract_inputs(self):
        inputs = []
        for input in self.sess.get_inputs():
            input_value = {
                'size': [modefied_dynamic(dim) for dim in input.shape],
                'name': input.name,
                'dtype':
                modified2np(modified2np(input.type[7:-1]))  # tensor(float)
            }
            inputs.append(input_value)
        return inputs

    def _extract_outputs(self):
        outputs = []
        for output in self.sess.get_outputs():
            output_value = {
                'size': [modefied_dynamic(dim) for dim in output.shape],
                'name': output.name,
                'dtype': modified2np(output.type[7:-1])  # tensor(float)
            }
            outputs.append(output_value)
        return outputs

    def _extract_ops(self):
        op_types = map(lambda x: x.op_type, self.nodes)
        ops = collections.Counter(op_types)
        return ops

    def _load_model(self):
        path = self._find_with_extension(EXTENSION)
        self.nodes = onnx.load_model(path).graph.node
        self.sess = onnxr.InferenceSession(path)
