# -*- coding: utf-8 -*-
import os
import sys
import shutil
import time
import json
from loguru import logger

from utils.get_model import check_model
from utils.config_generator import TRTISConfigGenerator
from utils.model_formatter import ModelFormatter

SKLEARN_MODEL = "model.joblib"
XGBOOST_MODEL = "model.xgboost"

class Preprocessor:
    """
    Preprocessor formats its directory structure and
    generate 'config.pbtxt'
    """
    env_list = [
        'MODEL_STORE', 'SERVING_NAME'
    ]

    def __init__(self):
        def get_critical_env(env):
            """
            A more strict way to get critical environment variable
            """
            try:
                var = os.environ[env]
            except KeyError as _:
                logger.error(f"{env} not defined")
                sys.exit(1)
            else:
                return var

        for env in Preprocessor.env_list:
            setattr(self, '_' + env.lower(), get_critical_env(env))

        self._trtis_conifig_generator = TRTISConfigGenerator()
        self.model_root_path = self._model_store
        self.model_path = os.path.join(self.model_root_path, self._serving_name, "1")
        

    def _extract_yaml(self):
        try:
            buffer, yaml_data = check_model(self._model_store, self._serving_name)
        except Exception as e:
            logger.error('error when checking model: ', e)
            sys.exit(1)
        else:
            logger.info(f'extract yaml_data succeed \n{buffer}')
            return yaml_data

    def _generate_config_pbtxt(self, yaml_data):
        try:
            config_pbtext_path = os.path.join(self.model_root_path, self._serving_name)
            self._trtis_conifig_generator.generate_config(yaml_data, config_pbtext_path, self._serving_name)
        except Exception as e:
            logger.error('error when generating config.pbtxt: ', e)
            sys.exit(1)
    
    def _generate_model_setting(self, format):
        setting = {}
        if format == 'SKLearn':
            setting = {
                'name': self._serving_name,
                'version': 'v0.1.0',
                'implementation': 'mlserver.models.SKLearnModel',
                'parameters': {
                    'uri': os.path.join(self.model_path, SKLEARN_MODEL)
                }
            }
        elif format == 'XGBoost':
            setting = {
                'name': self._serving_name,
                'version': 'v0.1.0',
                'implementation': 'mlserver.models.XGBoostModel',
                'parameters': {
                    'uri': os.path.join(self.model_path, XGBOOST_MODEL)
                }
            }

        json_str = json.dumps(setting)
        with open(os.path.join(self.model_root_path, "model-settings.json"), 'w') as json_file:
            json_file.write(json_str)

    def _format_model(self, format):
        format = format.lower()
        try:
            formatter = ModelFormatter(format)
            formatter.format(self.model_path)
        except Exception as e:
            logger.error(f'error when formatting directory {self.model_path}: ', e)
            sys.exit(1)

    def start(self):
        ormb_file_path = os.path.join(self.model_root_path, self._serving_name, "ormbfile.yaml")
        if not os.path.exists(ormb_file_path):
            return

        # Phase 1: Extract model_format and yaml
        yaml_data = self._extract_yaml()
        format = yaml_data["format"]

        # Phase 2: Generate 'config.pbtxt' if need
        if format != 'PMML' and format != 'SKLearn' and format != 'XGBoost':
           self._generate_config_pbtxt(yaml_data)

        # Phase 3: Generate 'model setting' if need
        if format == 'SKLearn' or format == 'XGBoost':
            self._generate_model_setting(format)

        # Phase 4: Re-organize directory format
        self._format_model(format)

        os.remove(ormb_file_path)  

if __name__ == '__main__':

    p = Preprocessor()
    p.start()
