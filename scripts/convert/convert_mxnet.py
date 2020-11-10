import argparse
import json
import os
import numpy as np
from mxnet.contrib import onnx as onnx_mxnet
from onnx import checker
import onnx
import sys

from base_convert.base_convert import BaseConvert

class MXNetToONNX(BaseConvert):
    def __init__(self):
        self.parser = argparse.ArgumentParser(
            description='Process some the path.')

        self.parser_args()

        self.args = self.parser.parse_args()
        self.input_dir = self.args.input_dir
        self.output_dir = self.args.output_dir
        self.model_name = self.args.model_name
        self.input_value = []
        try:            
            if os.getenv("MODEL_METADATA_FROM_ENV", "") != "":
                inputs = os.getenv("INPUTS", "")
                self.input_value = json.loads(inputs)
            else:
                ormbfile_content = self.read_ormbfile(self.input_dir)
                self.input_value = ormbfile_content["signature"]["inputs"]
        except Exception as e:
            raise e

        self.params_path = ""
        self.symbol_json = ""

    def parser_args(self):
        self.parser.add_argument(
            '--input_dir',
            help='path to the directory of model needs to be converted',
            dest='input_dir')
        self.parser.add_argument(
            '--output_dir',
            help='path to the directory of exported onnx model',
            dest='output_dir')
        self.parser.add_argument(
            '--model_name',
            help='the name of the model needs to be converted',
            dest='model_name',
            default="model.onnx")

    def _parase_modelfile(self):
        for root, _, files in os.walk(self.input_dir):
            for f in files:
                if ".params" in f:
                    self.params_path = os.path.join(root, f)
                if "symbol.json" in f:
                    self.symbol_json = os.path.join(root, f)
            if self.params_path != "" and self.symbol_json != "":
                break

        if self.params_path == "" or self.symbol_json == "":
            raise Exception("not found params_path or sybol_json")

    def _parse_input(self):
        self.input_shape = []
        self.data_type = None
        for input in self.input_value:
            self.input_shape.append(tuple(input['size']))
            if self.data_type is None:
                self.data_type = np.dtype(input['dtype'])

    def converter(self):
        self._parse_input()
        self._parase_modelfile()
        os.makedirs(self.output_dir, exist_ok=True)
        out_path = os.path.join(self.output_dir, 'model', self.model_name)
        converted_model_path = onnx_mxnet.export_model(self.symbol_json,
                                                       self.params_path,
                                                       self.input_shape,
                                                       self.data_type,
                                                       out_path)

        model_proto = onnx.load(converted_model_path)
        checker.check_graph(model_proto.graph)

        # fix mxnet bn bugs https://github.com/onnx/models/issues/156
        model = onnx.load(out_path)
        for node in model.graph.node:
            if (node.op_type == "BatchNormalization"):
                for attr in node.attribute:
                    if (attr.name == "spatial"):
                        attr.i = 1
        onnx.save(model, out_path)

        super().write_output_ormbfile(self.input_dir, self.output_dir)


if __name__ == '__main__':
    convert = MXNetToONNX()
    convert.converter()
