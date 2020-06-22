import os
import json
import collections

import mxnet as mx
from mxnet.gluon import nn

from .base_extract import BaseExtrctor

MIN_OPS = 3
MODEL_TYPE = "MXNet"
ExtenSymbol = 'symbol.json'
ExtenParams = '.params'


class MXNetExtractor(BaseExtrctor):
    def _extract_ops(self):
        def find_ops(module, l):
            if getattr(module, '_children'):
                for sub_module in module._children.values():
                    find_ops(sub_module, l)
            else:
                l.append(type(module).__name__)

        op_types = []
        find_ops(self.model, op_types)
        '''sometimes a huge module like a single op'''
        if len(op_types) < MIN_OPS:
            op_types = self.ops
        ops = collections.Counter(op_types)
        return ops

    def _load_model(self):
        def _isInput(s):
            return s.startswith('data') and (len(s) == 4 or str.isalnum(s[4:]))

        symbol_path = self._find_with_extension(ExtenSymbol)
        params_path = self._find_with_extension(ExtenParams)
        with open(symbol_path, 'rb') as f:
            tmp_json = json.load(f)
            self.ops = [node['op'] for node in tmp_json['nodes']]
            inputs_name = [
                node['name'] for node in tmp_json['nodes']
                if _isInput(node['name'])
            ]
        self.model = nn.SymbolBlock.imports(symbol_path,
                                            inputs_name,
                                            param_file=params_path)
