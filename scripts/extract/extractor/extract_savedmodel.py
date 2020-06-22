import os
import json
import collections

import tensorflow as tf
from google.protobuf import message
from google.protobuf import text_format
from tensorflow.core.protobuf import saved_model_pb2
from tensorflow.python.lib.io import file_io
from tensorflow.python.saved_model import constants
from tensorflow.python.util import compat

from .base_extract import BaseExtrctor

MODEL_TYPE = "SavedModel"
SIGNATURE_KEY = 'serving_default'


class TensorflowExtractor(BaseExtrctor):
    def _extract_inputs(self):
        inputs = []
        v = self.signature[SIGNATURE_KEY]
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
        v = self.signature[SIGNATURE_KEY]
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
        global SIGNATURE_KEY
        metaGraph = self.get_saved_model_tag_sets()
        self.graph = metaGraph.graph_def
        self.signature = metaGraph.signature_def
        assert len(
            self.signature) > 0, "No signature found, maybe its frozen graph"

        assert (SIGNATURE_KEY in self.signature) or (len(
            self.signature
        ) == 1), "%s must in signature or only one signature in savedmodel" % (
            SIGNATURE_KEY)

        if not SIGNATURE_KEY in self.signature:
            SIGNATURE_KEY = next(iter(self.signature.keys()))

    def get_saved_model_tag_sets(self):
        """Retrieves all the tag-sets available in the SavedModel.
        Args:
          saved_model_dir: Directory containing the SavedModel.
        Returns:
          String representation of all tag-sets in the SavedModel.
        """
        saved_model = self.read_saved_model()
        assert len(saved_model.meta_graphs) >= 1, 'no metagraph found'
        return saved_model.meta_graphs[0]

    def read_saved_model(self):
        """Reads the savedmodel.pb or savedmodel.pbtxt file containing `SavedModel`.
        Args:
          saved_model_dir: Directory containing the SavedModel file.
        Returns:
          A `SavedModel` protocol buffer.
        Raises:
          IOError: If the file does not exist, or cannot be successfully parsed.
        """
        # Build the path to the SavedModel in pbtxt format.
        path_to_pbtxt = os.path.join(
            compat.as_bytes(self.dir),
            compat.as_bytes("model"),
            compat.as_bytes(constants.SAVED_MODEL_FILENAME_PBTXT))
        # Build the path to the SavedModel in pb format.
        path_to_pb = os.path.join(
            compat.as_bytes(self.dir),
            compat.as_bytes("model"),
            compat.as_bytes(constants.SAVED_MODEL_FILENAME_PB))

        # Ensure that the SavedModel exists at either path.
        if not file_io.file_exists(path_to_pbtxt) and not file_io.file_exists(
                path_to_pb):
            raise IOError("SavedModel file does not exist at: %s" % self.dir)

        # Parse the SavedModel protocol buffer.
        saved_model = saved_model_pb2.SavedModel()
        if file_io.file_exists(path_to_pb):
            try:
                file_content = file_io.FileIO(path_to_pb, "rb").read()
                saved_model.ParseFromString(file_content)
                return saved_model
            except message.DecodeError as e:
                raise IOError("Cannot parse file %s: %s." %
                              (path_to_pb, str(e)))
        elif file_io.file_exists(path_to_pbtxt):
            try:
                file_content = file_io.FileIO(path_to_pbtxt, "rb").read()
                text_format.Merge(file_content.decode("utf-8"), saved_model)
                return saved_model
            except text_format.ParseError as e:
                raise IOError("Cannot parse file %s: %s." %
                              (path_to_pbtxt, str(e)))
        else:
            raise IOError("SavedModel file does not exist at: %s/{%s|%s}" %
                          (self.dir, constants.SAVED_MODEL_FILENAME_PBTXT,
                           constants.SAVED_MODEL_FILENAME_PB))
