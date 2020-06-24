import argparse
import os
import tensorflow.keras as tf_keras
import tensorflow as tf
import keras
import keras.backend as K

from base_convert.base_convert import BaseConvert

class KerasToTensorFlow(BaseConvert):
    def __init__(self):
        self.parser = argparse.ArgumentParser(
            description='Process some the path.')

        self.parser_args()

        self.args = self.parser.parse_args()
        self.input_dir = self.args.input_dir
        self.ormbfile_output_dir = self.args.output_dir
        self.output_dir = os.path.join(self.args.output_dir, 'model')
        self.model_name = self.args.model_name

        self.model_path = os.path.join(self.input_dir, 'model', 'model.h5')

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
            default='model.h5',
            help='the name of the model needs to be converted',
            dest='model_name')

    def converter(self):
        try:
            tf.reset_default_graph()
            model = tf_keras.models.load_model(self.model_path)
            tf.contrib.saved_model.save_keras_model(model, self.output_dir)
            os.system('mv  %s/*/* %s' % (self.output_dir, self.output_dir))
            super().write_output_ormbfile(self.input_dir, self.ormbfile_output_dir)
        except Exception as ea:
            try:
                tf.reset_default_graph()
                sess = tf.Session()
                K.set_session(sess)
                model = keras.models.load_model(self.model_path)
                input_tensor = {}
                output_tensor = {}
                ## collect input nodes
                for i, i_input in enumerate(model.inputs):
                    input_tensor['input_%d' % i] = i_input
                ## collect output nodes
                for i, i_output in enumerate(model.outputs):
                    output_tensor['ouput_%d' % i] = i_output
                tf.saved_model.simple_save(sess,
                                           export_dir=self.output_dir,
                                           inputs=input_tensor,
                                           outputs=output_tensor)
                super().write_output_ormbfile(self.input_dir, self.ormbfile_output_dir)
            except Exception as eb:
                raise eb


if __name__ == '__main__':
    convert = KerasToTensorFlow()
    convert.converter()
