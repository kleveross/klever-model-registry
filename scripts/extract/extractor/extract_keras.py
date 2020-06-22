import os
import json
import collections

from tensorflow import keras

from .base_extract import BaseExtrctor

MODEL_TYPE = 'Keras'
EXTENSION = '.h5'


class KerasExtractor(BaseExtrctor):
    def _extract_inputs(self):
        inputs = []
        for tensor in self.models.inputs:
            origin_inputs = {}
            origin_inputs['name'] = tensor.name.split(':')[0]
            origin_inputs['dataType'] = tensor.dtype.as_numpy_dtype.__name__
            origin_inputs['dims'] = [
                i.value if i.value else -1 for i in tensor.shape
            ]
            inputs.append(origin_inputs)
        return inputs

    def _extract_outputs(self):
        outputs = []
        for tensor in self.models.outputs:
            origin_outputs = {}
            origin_outputs['name'] = tensor.name.split(':')[0]
            origin_outputs['dataType'] = tensor.dtype.as_numpy_dtype.__name__
            origin_outputs['dims'] = [
                i.value if i.value else -1 for i in tensor.shape
            ]
            outputs.append(origin_outputs)
        return outputs

    def _extract_ops(self):
        origin_ops = [node.op for node in self.graph.node]
        ops = collections.Counter(origin_ops)
        return ops

    def _load_model(self):
        path = self._find_with_extension(EXTENSION)
        try:
            self.models = keras.models.load_model(path)
        except Exception as e:
            raise IOError("Cannot read file %s: %s." % (path, str(e)))

        with keras.backend.get_session() as sess:
            self.graph = sess.graph_def
