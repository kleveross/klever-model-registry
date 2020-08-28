# -*- coding: utf-8 -*-
import os
import yaml
from loguru import logger

import utils.help_functions as hp


@logger.catch()
class TRTISConfigGenerator(object):
    """
    TRTISConfigGenerator generates 'config.pbtxt' for
    TensorRT Inference Server (TRTIS)
    """
    _implemented_runtimes = [
        'onnxruntime_onnx', 'pmmlruntime_pmml', 'pytorch_libtorch', 'tensorflow_savedmodel',
        'tensorflow_graphdef', 'caffe2_netdef', 'tensorrt_plan'
    ]

    _dataType_dict = {
        'bool': 'TYPE_BOOL',
        'uint8': 'TYPE_UINT8',
        'uint16': 'TYPE_UINT16',
        'uint64': 'TYPE_UINT64',
        'int8': 'TYPE_INT8',
        'int16': 'TYPE_INT16',
        'int32': 'TYPE_INT32',
        'int64': 'TYPE_INT64',
        'float16': 'TYPE_FP16',
        'float32': 'TYPE_FP32',
        'float64': 'TYPE_FP64',
        'object': 'TYPE_STRING'
    }

    def __init__(self):
        self.overall_template = \
            '''
            name: "{name}"
            platform: "{platform}"
            {max_batch_size_content}
            input [
            {all_inputs}
            ]
            output [
            {all_outputs}
            ]
            '''
        self.xput_template = \
            '''
            {{
                name: "{name}"
                data_type: {dataType}
                dims: {dims}
            }}
            '''

    @logger.catch()
    def _fixed_dim(self, xputs):
        if not isinstance(xputs, list):
            xputs = [xputs]

        for i in xputs:
            dims = i['dims']
            if isinstance(dims, list):
                if len(dims) > 0:
                    if dims[0] == -1:
                        return False

        return True

    @logger.catch()
    def _gen_xputs(self, xputs):
        if not isinstance(xputs, list):
            xputs = [xputs]

        # Assume there is no duplications
        xputs_str_list = [
            self.xput_template.format(
                name=xput['name'],
                dataType=TRTISConfigGenerator._dataType_dict[xput['dtype']],
                dims=str(xput['size'])) for xput in xputs
        ]

        return ',\n'.join(xputs_str_list)

    @logger.catch()
    def generate_config(self, manifest, dst_dir, serving_name):
        """
        generate_config: generate the content of 'config.pbtxt' and write it to disk

        Args:
            scene_name: str, the name of scene where the serving is deployed
            manifest: json structure, read content from 'ormbfile.yaml'
            dst_dir: str, location where the 'config.pbtxt' should be stored

        Return: None
        """
        max_bs = 1
        format = manifest['format'].lower()
        platform = hp.get_platform_by_format(format)
        assert platform in TRTISConfigGenerator._implemented_runtimes

        inputs_str = self._gen_xputs(manifest['signature']['inputs'])
        outputs_str = self._gen_xputs(manifest['signature']['outputs'])

        max_batch_size_content = '' #('' if self._fixed_dim(manifest['inputs'])
                                 # else f"max_batch_size: {max_bs}")

        config_pbtxt = self.overall_template.format(
            name=serving_name,
            platform=platform,
            max_batch_size_content=max_batch_size_content,
            all_inputs=inputs_str,
            all_outputs=outputs_str)
        logger.info(config_pbtxt)
        with open(os.path.join(dst_dir, 'config.pbtxt'), 'w') as f:
            f.write(config_pbtxt)


# For Testing Only
# Please set the following 3 environment variables
# 1. TEST_YAML_PATH: the path to the 'ormbfile.yaml'
# 2. TEST_SCENE_NAME: any scene name for test
# 3. TEST_TEMP_DIR: the destination directory for 'config.pbtxt'
# For example:
# TEST_YAML_PATH=/tmp/test_onnx/raw/ormbfile.yaml TEST_SCENE_NAME=haha_test \
# TEST_TEMP_DIR=/tmp/test_onnx python3 config_generator.py
if __name__ == '__main__':
    manifest_read = yaml.load(open(os.environ['TEST_YAML_PATH']).read())
    generator = TRTISConfigGenerator()
    generator.generate_config(os.environ['TEST_SCENE_NAME'], manifest_read,
                              os.environ['TEST_TEMP_DIR'], "models")
