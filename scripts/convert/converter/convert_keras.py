import argparse
import os
import tensorflow.keras as tf_keras
import tensorflow as tf

from .base_convert import BaseConverter

EXTENSTION = '.h5'


class KerasToTensorFlow(BaseConverter):
    def _load_model(self):
        self.model_path = self._find_with_extension(EXTENSTION)

    def _convert(self):
        try:
            tf.reset_default_graph()
            model = tf_keras.models.load_model(self.model_path)
            out_path = os.path.join(self.output_dir, 'model')
            model.save(out_path, save_format='tf')
        except Exception as ea:
            raise ea


if __name__ == '__main__':
    convert = KerasToTensorFlow()
    convert.convert()
