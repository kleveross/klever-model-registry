# -*- coding: utf-8 -*-
import os
import shutil

import utils.help_functions as hp


class ModelFormatInterface:
    def __init__(self):
        pass

    def execute(self):
        raise NotImplementedError


class ONNXFormatter(ModelFormatInterface):
    _target_onnx_filename = 'model.onnx'

    def execute(self, target_dir):
        onnx_file = hp.find_file_ends_with(target_dir, '.onnx')
        assert len(onnx_file) == 1
        hp.rename(target_dir, onnx_file[0],
                  ONNXFormatter._target_onnx_filename)


class TFSavedModelFormatter(ModelFormatInterface):
    _target_savedmodel_dirname = 'model.savedmodel'

    def execute(self, target_dir):
        # 1. mkdir ./model.savedmodel
        os.makedirs(
            os.path.join(target_dir,
                         TFSavedModelFormatter._target_savedmodel_dirname))
        # 2. find . | grep "[.]pb"
        pb_file = hp.find_file_ends_with(target_dir, '.pb')
        assert len(pb_file) == 1
        hp.rename(
            target_dir, pb_file[0],
            os.path.join(TFSavedModelFormatter._target_savedmodel_dirname,
                         pb_file[0]))
        if len(hp.find_file_ends_with(target_dir, 'variables')) > 0:
            shutil.move(
                os.path.join(target_dir, 'variables'),
                os.path.join(target_dir,
                             TFSavedModelFormatter._target_savedmodel_dirname))


class TorchScriptFormatter(ModelFormatInterface):
    _target_torchscript_filename = 'model.pt'

    def execute(self, target_dir):
        pt_file = hp.find_file_ends_with(target_dir, '.pt')
        assert len(pt_file) == 1
        hp.rename(target_dir, pt_file[0],
                  TorchScriptFormatter._target_torchscript_filename)


class Caffe2NetDefFormatter(ModelFormatInterface):
    _target_model_filename = 'model.netdef'
    _target_init_model_filename = 'init_model.netdef'

    def execute(self, target_dir):
        model_file = hp.find_file_ends_with(target_dir, 'predict_net.pb')
        init_file = hp.find_file_ends_with(target_dir, 'init_net.pb')

        assert len(model_file) == 1 and len(init_file) == 1
        hp.rename(target_dir, model_file[0],
                  Caffe2NetDefFormatter._target_model_filename)
        hp.rename(target_dir, init_file[0],
                  Caffe2NetDefFormatter._target_init_model_filename)


class TFGraphDefFormatter(ModelFormatInterface):
    _target_graphdef_filename = 'model.graphdef'

    def execute(self, target_dir):
        pb_file = hp.find_file_ends_with(target_dir, '.pb')
        assert len(pb_file) == 1
        hp.rename(target_dir, pb_file[0],
                  TFGraphDefFormatter._target_graphdef_filename)


class PMMLFormatter(ModelFormatInterface):
    _target_pmml_filename = 'model.pmml'

    def execute(self, target_dir):
        pmml_file = hp.find_file_ends_with(target_dir, '.pmml')
        assert len(pmml_file) == 1
        hp.rename(target_dir, pmml_file[0],
                  PMMLFormatter._target_pmml_filename)


class TensorRTFormatter(ModelFormatInterface):
    _target_plan_filename = 'model.plan'

    def execute(self, target_dir):
        tensorrt_file = hp.find_file_ends_with(target_dir, '.plan')
        assert len(tensorrt_file) == 1
        hp.rename(target_dir, tensorrt_file[0],
                  TensorRTFormatter._target_plan_filename)


class SKLearnFormatter(ModelFormatInterface):
    _target_sklearn_filename = 'model.joblib'

    def execute(self, target_dir):
        sklearn_file = hp.find_file_ends_with(target_dir, '.joblib')
        assert len(sklearn_file) == 1
        hp.rename(target_dir, sklearn_file[0],
                  SKLearnFormatter._target_sklearn_filename)


class XGBoostFormatter(ModelFormatInterface):
    _target_xgboost_filename = 'model.xgboost'

    def execute(self, target_dir):
        xgboost_file = hp.find_file_ends_with(target_dir, '.xgboost')
        assert len(xgboost_file) == 1
        hp.rename(target_dir, xgboost_file[0],
                  XGBoostFormatter._target_xgboost_filename)


class ModelFormatter:
    _implemented_dict = {
        'onnxruntime_onnx': ONNXFormatter,
        'tensorflow_savedmodel': TFSavedModelFormatter,
        'pytorch_libtorch': TorchScriptFormatter,
        'tensorflow_graphdef': TFGraphDefFormatter,
        'caffe2_netdef': Caffe2NetDefFormatter,
        'pmmlruntime_pmml': PMMLFormatter,
        'tensorrt_plan': TensorRTFormatter,
        'scikitlearn_sklearn': SKLearnFormatter,
        'xgboost_xgboost': XGBoostFormatter
    }

    def __init__(self, format):
        platform = hp.get_platform_by_format(format)
        assert platform in ModelFormatter._implemented_dict
        self._formatter = ModelFormatter._implemented_dict[platform]()

    def format(self, target_dir):
        self._formatter.execute(target_dir)
