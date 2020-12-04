import os
import json
import collections

import tensorflow as tf

from .base_extract import BaseExtrctor

MODEL_TYPE = 'GraphDef'
EXTENSION = '.graphdef'


class TensorflowExtractor(BaseExtrctor):
    def _extract_inputs(self):
        inputs = []
        for node in self.graph_def.node:
            if node.op == 'Placeholder':
                origin_inputs = {}
                origin_inputs['name'] = node.name
                origin_inputs['dType'] = tf.DType(
                    node.attr['dType'].type).as_numpy_dtype.__name__
                dim = node.attr['shape'].shape.dim
                origin_inputs['size'] = [dim[i].size for i in range(len(dim))]
                inputs.append(origin_inputs)
        return inputs

    def _extract_outputs(self):
        outputs = []
        all_set = set()
        inner_set = set()
        for node in self.graph_def.node:
            all_set.add(node.name)
            for father in node.input:
                inner_set.add(father)

        output_set = all_set - inner_set
        outputs_name = [i + ':0' for i in output_set]

        outputs_tensor = tf.import_graph_def(self.graph_def,
                                             name='',
                                             return_elements=outputs_name)
        for tensor in outputs_tensor:
            origin_outputs = {}
            origin_outputs['name'] = tensor.name.split(':')[0]
            origin_outputs['dType'] = tensor.dtype.as_numpy_dtype.__name__
            origin_outputs['size'] = [
                i.value if i.value else -1 for i in tensor.shape
            ]
            outputs.append(origin_outputs)
        return outputs

    def _extract_ops(self):
        origin_ops = [node.op for node in self.graph_def.node]
        ops = collections.Counter(origin_ops)
        return ops

    def _load_model(self):
        path = self._find_with_extension(EXTENSION)
        with tf.gfile.FastGFile(path, 'rb') as f:
            graph_def = tf.GraphDef()
            graph_def.ParseFromString(f.read())
        self.graph_def = graph_def
