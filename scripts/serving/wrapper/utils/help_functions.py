# -*- coding: utf-8 -*-
import os
from loguru import logger

mlserver_model = [
    'SKLearn', 'XGBoost', 'MLlib'
]


@logger.catch()
def find_file_ends_with(dir_path, ext):
    ext_upper = ext.upper()
    ext_lower = ext.lower()
    ojbect_list = os.listdir(dir_path)

    return list(
        filter(lambda x: x.endswith(ext_lower) or x.endswith(ext_upper),
               ojbect_list))


@logger.catch()
def rename(dir_path, ori_name, new_name):
    ori_name_full = os.path.join(dir_path, ori_name)
    new_name_full = os.path.join(dir_path, new_name)
    os.rename(ori_name_full, new_name_full)


@logger.catch()
def isTritonModel(format):
    if format != 'PMML' and format not in mlserver_model:
        return True
    return False


@logger.catch()
def isMLServerModel(format):
    if format in mlserver_model:
        return True
    return False


@logger.catch()
def get_platform_by_format(format):
    format_platform_dict = {
        'onnx': 'onnxruntime_onnx',
        'savedmodel': 'tensorflow_savedmodel',
        'torchscript': 'pytorch_libtorch',
        'graphdef': 'tensorflow_graphdef',
        'netdef': 'caffe2_netdef',
        'pmml': 'pmmlruntime_pmml',
        'tensorrt': 'tensorrt_plan',
        'sklearn': 'scikitlearn_sklearn',
        'xgboost': 'xgboost_xgboost',
        'mllib': 'mllib_mllib'
    }

    return format_platform_dict[format]
