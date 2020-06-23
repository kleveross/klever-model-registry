import os
import json
import collections

import onnx
import onnx.helper
from .base_extract import BaseExtrctor

MODEL_TYPE = 'ONNX'
EXTENSION = '.onnx'


def modified2np(s):
    if s == 'float':
        return 'float32'
    if s == 'double':
        return 'float64'
    return s


modefied_dynamic = lambda x: -1 if not x else x


class OnnxExtractor(BaseExtrctor):
    def _extract_inputs(self):
        origin_inputs = list(
            filter(lambda x: x.name not in self.initialized,
                   self.model.graph.input))
        inputs = []
        for input in origin_inputs:
            input_value = {
                'size': [
                    modefied_dynamic(dim.dim_value)
                    for dim in input.type.tensor_type.shape.dim
                ],
                'name':
                input.name,
                'dtype':
                modified2np(onnx.TensorProto.DataType.keys()[
                    input.type.tensor_type.elem_type].lower())
            }
            inputs.append(input_value)
        return inputs

    def _extract_outputs(self):
        origin_outputs = list(self.model.graph.output)
        outputs = []
        for output in origin_outputs:
            output_value = {
                'size': [
                    modefied_dynamic(dim.dim_value)
                    for dim in output.type.tensor_type.shape.dim
                ],
                'name':
                output.name,
                'dtype':
                modified2np(onnx.TensorProto.DataType.keys()[
                    output.type.tensor_type.elem_type].lower())
            }
            outputs.append(output_value)
        return outputs

    def _extract_ops(self):
        op_types = map(lambda x: x.op_type, self.model.graph.node)
        ops = collections.Counter(op_types)
        return ops

    def _load_model(self):
        path = self._find_with_extension(EXTENSION)
        self.model = onnx.load_model(path)
        self.initialized = {t.name for t in self.model.graph.initializer}
