import argparse
import json
import os
import numpy as np
from mxnet.contrib import onnx as onnx_mxnet
from onnx import checker
import onnx

from base_convert import BaseConverter

ExtenSymbol = 'symbol.json'
ExtenParams = '.params'
MODEL_NAME = 'model.onnx'


class MXNetToONNX(BaseConverter):
    def _load_model(self):
        self.params_path = self._find_with_extension(ExtenParams)
        self.symbol_json = self._find_with_extension(ExtenSymbol)

    def _parse_input(self):
        self.input_shape = []
        self.data_type = None
        for input in self.input_value:
            self.input_shape.append(tuple(input['Dims']))
            if self.data_type is None:
                self.data_type = np.dtype(input['DataType'])

    def _convert(self):
        self._load_model()
        self._parse_input()
        out_path = os.path.join(self.output_dir, 'model', MODEL_NAME)
        converted_model_path = onnx_mxnet.export_model(self.symbol_json,
                                                       self.params_path,
                                                       self.input_shape,
                                                       self.data_type,
                                                       out_path)

        model_proto = onnx.load(converted_model_path)
        checker.check_graph(model_proto.graph)


if __name__ == '__main__':
    convert = MXNetToONNX()
    convert.convert()
