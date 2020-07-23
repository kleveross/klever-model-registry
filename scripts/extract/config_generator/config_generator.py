import os
import json


class TritonConfigGenerator(object):
    def __init__(self):
        self.overall_template = \
'''name: "{alias}"
platform: "{platform}"
max_batch_size: "{max_bs}"
input [
{all_inputs}
]
output [
{all_outputs}
]
'''
        self.xput_template = \
'''  {{ 
    name: "{xput_name}"
    data_type: {dtype}
    dims: {dims}
  }}'''

        # Referring: https://docs.nvidia.com/deeplearning/triton-inference-server/master-user-guide/docs/model_configuration.html#datatypes
        # as well as https://github.com/oracle/graphpipe/blob/master/docs/guide/user-guide/spec.md
        self.type_dict = {
            'bool': 'TYPE_BOOL',
            'uint8': 'TYPE_UINT8',
            'uint16': 'TYPE_UINT16',
            'uint32': 'TYPE_UINT32',
            'uint64': 'TYPE_UINT64',
            'int8': 'TYPE_INT8',
            'int16': 'TYPE_INT16',
            'int32': 'TYPE_INT32',
            'int64': 'TYPE_INT64',
            'float16': 'TYPEFP16',
            'float32': 'TYPE_FP32',
            'float64': 'TYPE_FP64',
            'string': 'TYPE_STRING'
        }

    def _gen_xputs(self, xputs):
        if not isinstance(xputs, list):
            xputs = [xputs]
        xputs_str = list(
            map(
                lambda xput: self.xput_template.format(
                    xput_name=xput['name'],
                    dtype=self.type_dict[str(xput['dtype'])],
                    dims=str(xput['size'][:])), xputs))
        return ',\n'.join(xputs_str)

    def generate_config(self, ormbfile, config_dir):
        signature = ormbfile['signature']
        inputs_str  = self._gen_xputs(signature['inputs']) 
        outputs_str = self._gen_xputs(signature['outputs'])
        if os.environ['FRAMEWORK'] == 'TensorFlow':
            max_bs = 0
        else:
            max_bs = 1

        if len(os.environ['FRAMEWORK']) != 0 and len(os.environ['FORMAT']) != 0:
            platform = '{}_{}'.format(os.environ['FRAMEWORK'].lower(), os.environ['FORMAT'].lower())
        else:
            platform ='custom'

        _, alias = os.path.split(os.environ['SOURCE_MODEL_TAG'])

        res = self.overall_template.format(
            alias=alias,
            platform=platform,
            max_bs=max_bs,
            all_inputs=inputs_str,
            all_outputs=outputs_str)
        print(res)
        with open(os.path.join(config_dir, 'model', 'config.pbtxt'), 'w') as f:
            f.write(res)
