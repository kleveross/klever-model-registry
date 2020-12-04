import os
import json
import collections

import re
from torch import jit

from .base_extract import BaseExtrctor

MODEL_TYPE = 'TorchScript'
EXTENSION = '.pt'


class TorchscriptExtractor(BaseExtrctor):
    def _extract_ops(self):
        get_net = lambda x: x._modules.items()
        fetch_f = lambda x: re.findall('\D+', x)[0]

        def find_ops(module, l, k=''):
            if len(get_net(module)) == 0:
                l.append(fetch_f(k))
            else:
                for k, v in get_net(module):
                    find_ops(v, l, k)

        op_types = []
        find_ops(self.model, op_types)
        layers = collections.Counter(op_types)
        return layers

    def _load_model(self):
        PATH = self._find_with_extension(EXTENSION)
        self.model = jit.load(PATH)
