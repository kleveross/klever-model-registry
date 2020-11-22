import os
import json
import collections

import tensorflow as tf
from google.protobuf import message
from google.protobuf import text_format
from tensorflow.core.protobuf import saved_model_pb2
from tensorflow.python.lib.io import file_io
import tensorflow.saved_model as saved_model
from tensorflow.python.util import compat

from .base_extract import BaseExtrctor

MODEL_TYPE = 'SavedModel'
DEFAULT_SERVING_SIGNATURE_DEF_KEY = 'serving_default'
INIT_OP_SIGNATURE_DEF_KEY = '__saved_model_init_op'
TRAIN_OP_SIGNATURE_DEF_KEY = '__saved_model_train_op'


class TensorflowExtractor(BaseExtrctor):
    def _extract_inputs(self):
        inputs = []
        v = self.signature
        for key in v.inputs:
            origin_inputs = {}
            origin_inputs['signatureConst'] = key
            origin_inputs['name'] = v.inputs[key].name
            origin_inputs['dtype'] = tf.DType(
                v.inputs[key].dtype).as_numpy_dtype.__name__
            dim = v.inputs[key].tensor_shape.dim
            origin_inputs['size'] = [dim[i].size for i in range(len(dim))]
            inputs.append(origin_inputs)
        return inputs

    def _extract_outputs(self):
        outputs = []
        v = self.signature
        for key in v.outputs:
            origin_outputs = {}
            origin_outputs['signatureConst'] = key
            origin_outputs['name'] = v.outputs[key].name
            origin_outputs['dtype'] = tf.DType(
                v.outputs[key].dtype).as_numpy_dtype.__name__
            dim = v.outputs[key].tensor_shape.dim
            origin_outputs['size'] = [dim[i].size for i in range(len(dim))]
            outputs.append(origin_outputs)
        return outputs

    def _extract_ops(self):
        origin_ops = [node.op for node in self.graph.node]
        ops = collections.Counter(origin_ops)
        return ops

    def _load_model(self):
        """Retrieves all the tag-sets available in the SavedModel.
        Args:
          saved_model_dir: Directory containing the SavedModel.
        Returns:
          String representation of all tag-sets in the SavedModel.
        """
        sess = tf.Session()
        MetaGraphDef = saved_model.load(sess, [saved_model.SERVING], self.dir)
        sig = None
        if DEFAULT_SERVING_SIGNATURE_DEF_KEY not in MetaGraphDef.signature_def:
            for sig_itr in MetaGraphDef.signature_def:
                if (sig_itr != INIT_OP_SIGNATURE_DEF_KEY
                        and sig_itr != TRAIN_OP_SIGNATURE_DEF_KEY):
                    sig = sig_itr
                    break
        else:
            sig = DEFAULT_SERVING_SIGNATURE_DEF_KEY

        assert sig is not None, "unable to load model , expected " + DEFAULT_SERVING_SIGNATURE_DEF_KEY + " signature"

        self.graph = MetaGraphDef.graph_def
        self.signature = MetaGraphDef.signature_def[sig]
