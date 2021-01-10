# -*- coding: utf-8 -*-
import os
import sys
import json
import shutil
from loguru import logger

from utils.get_model import check_model
from utils.config_generator import TRTISConfigGenerator
from utils.model_formatter import ModelFormatter
from utils.help_functions import isTritonModel, isMLServerModel

SKLEARN_MODEL = 'model.joblib'
XGBOOST_MODEL = 'model.xgboost'

USING_ORMBFILE_ENV = "USING_ORMBFILE"
INPUTS_ENV = "INPUTS"
OUTPUTS_ENV = "OUTPUTS"
FORMAT_ENV = "FORMAT"

class Preprocessor:
    '''
    Preprocessor formats its directory structure and
    generate 'config.pbtxt'
    '''
    env_list = [
        'MODEL_STORE', 'SERVING_NAME'
    ]

    def __init__(self):
        def get_critical_env(env):
            '''
            A more strict way to get critical environment variable
            '''
            try:
                var = os.environ[env]
            except KeyError:
                logger.error(f'{env} not defined')
                sys.exit(1)
            else:
                return var

        for env in Preprocessor.env_list:
            setattr(self, '_' + env.lower(), get_critical_env(env))
        
        self.using_ormbfile = True if os.getenv(USING_ORMBFILE_ENV, '') else False
        self._trtis_conifig_generator = TRTISConfigGenerator()
        self.model_root_path = self._model_store
        self.model_path = os.path.join(
            self.model_root_path, self._serving_name, '1')

    def _extract_yaml(self):
        try:
            buffer, yaml_data = check_model(
                self._model_store, self._serving_name)
        except Exception as e:
            logger.error('error when checking model: ', e)
            sys.exit(1)
        else:
            logger.info(f'extract yaml_data succeed \n{buffer}')
            return yaml_data

    def _generate_config_pbtxt(self, yaml_data):
        try:
            config_pbtext_path = os.path.join(
                self.model_root_path, self._serving_name)
            self._trtis_conifig_generator.generate_config(
                yaml_data, config_pbtext_path, self._serving_name)
        except Exception as e:
            logger.error('error when generating config.pbtxt: ', e)
            sys.exit(1)

    def _generate_model_setting(self, format, version):
        setting = {}
        if format == 'SKLearn':
            setting = {
                'name': self._serving_name,
                'version': version,
                'implementation': 'mlserver.models.SKLearnModel',
                'parameters': {
                    'uri': os.path.join(self.model_path, SKLEARN_MODEL)
                }
            }
        elif format == 'XGBoost':
            setting = {
                'name': self._serving_name,
                'version': version,
                'implementation': 'mlserver.models.XGBoostModel',
                'parameters': {
                    'uri': os.path.join(self.model_path, XGBOOST_MODEL)
                }
            }
        elif format == 'MLlib':
            try:
                mllibformat = os.environ['MLLIB_FORMAT']
            except KeyError:
                logger.error('MLLIB_FORMAT not defined')
                sys.exit(1)

            setting = {
                'name': self._serving_name,
                'version': version,
                'implementation': 'mlservermllib.models.MLLibModel',
                'parameters': {
                    'uri': os.path.join(
                        self.model_root_path, self._serving_name, '1'),
                    'format': mllibformat
                }
            }

        json_str = json.dumps(setting)
        with open(os.path.join(self.model_root_path, 'model-settings.json'), 'w') as json_file:
            json_file.write(json_str)

    def _format_model(self, format):
        format = format.lower()
        try:
            formatter = ModelFormatter(format)
            formatter.format(self.model_path)
        except Exception as e:
            logger.error(
                f'error when formatting directory {self.model_path}: ', e)
            sys.exit(1)

    def start(self):
        ormb_file_path = os.path.join(
            self.model_root_path, self._serving_name, 'ormbfile.yaml')
        if not os.path.exists(ormb_file_path):
            logger.warn(f'{ormb_file_path} does not exist')
            return

        old_model_path = os.path.join(self.model_root_path, self._serving_name, 'model')
        if not os.path.isdir(old_model_path):
            logger.warn(f'{old_model_path} does not exist')
            return

        # Phase 1: Extract model_format and yaml
        format = ""
        yaml_data = {}
        if self.using_ormbfile:
            yaml_data = self._extract_yaml()
            if 'format' not in yaml_data.items():
                logger.error('model format missing')
                return
            format = yaml_data["format"]
        else:
            inputs = os.getenv(INPUTS_ENV, "[]")
            outputs = os.getenv(OUTPUTS_ENV, "[]")
            format = os.getenv(FORMAT_ENV, "")
            yaml_data["format"] = format
            yaml_data["signature"] = {}
            yaml_data["signature"]["inputs"] = json.loads(inputs)
            yaml_data["signature"]["outputs"] = json.loads(outputs)

        # Phase 2: Generate 'config.pbtxt' for triton models
        if isTritonModel(format):
            self._generate_config_pbtxt(yaml_data)

        # Phase 3: Generate 'model setting' for mlserver models
        if isMLServerModel(format):
            # set env for mlserver
            os.putenv('MODEL_FORMAT', format)

            # get version from ormbfile
            if 'version' in yaml_data.items():
                version = yaml_data['version']
            else:
                version = 'v1.0.0'

            self._generate_model_setting(format, version)

        # Phase 4: Re-organize directory format
        if os.path.isdir(self.model_path):
            shutil.rmtree(self.model_path)
        os.rename(old_model_path, self.model_path)
        self._format_model(format)

if __name__ == '__main__':
    p = Preprocessor()
    p.start()
