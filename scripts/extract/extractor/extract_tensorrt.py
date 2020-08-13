import os
import json
import collections

import tensorrt as trt

from .base_extract import BaseExtrctor

MODEL_TYPE = 'TensorRT'
MODEL_EXTENSION1 = '.engine'
MODEL_EXTENSION2 = '.plan'


class TensorrtExtractor(BaseExtrctor):
    def _extract_inputs(self):
        engine = self.engine
        inputs = []
        for binding in engine:
            if engine.binding_is_input(binding):
                origin_inputs = {}
                origin_inputs['name'] = binding
                origin_inputs['dtype'] = trt.nptype(
                    engine.get_binding_dtype(binding)).__name__
                origin_inputs['size'] = engine.get_binding_shape(binding)
                inputs.append(origin_inputs)
        return inputs

    def _extract_outputs(self):
        engine = self.engine
        outputs = []
        for binding in engine:
            if not engine.binding_is_input(binding):
                origin_outputs = {}
                origin_outputs['name'] = binding
                origin_outputs['dtype'] = trt.nptype(
                    engine.get_binding_dtype(binding)).__name__
                origin_outputs['size'] = engine.get_binding_shape(binding)
                outputs.append(origin_outputs)
        return outputs

    def _extract_ops(self):
        return {'layers': self.engine.num_layers}

    def _load_model(self):
        try:
            model_path = self._find_with_extension(MODEL_EXTENSION1)
        except:
            model_path = self._find_with_extension(MODEL_EXTENSION2)
        TRT_LOGGER = trt.Logger()
        with open(model_path, "rb") as f, trt.Runtime(TRT_LOGGER) as runtime:
            self.engine = runtime.deserialize_cuda_engine(f.read())