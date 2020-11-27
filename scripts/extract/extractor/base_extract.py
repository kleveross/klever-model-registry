import os
import json

MODEL_TYPE = 'BASE'


class BaseExtrctor(object):
    def __init__(self, path):
        self.dir = path

    def _extract_inputs(self):
        return []

    def _extract_outputs(self):
        return []

    def _extract_ops(self):
        raise NotImplementedError

    def extract(self):
        self._load_model()
        inputs = self._extract_inputs()
        outputs = self._extract_outputs()
        ops = self._extract_ops()
        res_dict = {
            "Inputs": inputs,
            "Outputs": outputs,
            "Operators": ops,
        }
        return res_dict

    def _find_with_extension(self, extension):
        dir = os.path.join(self.dir, "model")
        filelist = list(
            filter(lambda f: f.endswith(extension) and not f.startswith('.'),
                   os.listdir(dir)))
        assert (len(filelist) == 1), "expected one %s file,but found %s" % (
            extension, len(filelist))
        return os.path.join(dir, filelist[0])

    def _find_with_name(self, name):
        dir = os.path.join(self.dir, "model")
        filelist = list(filter(lambda f: f == name, os.listdir(dir)))
        assert (len(filelist) == 1
                ), "expected one %s file,but found %s" % (name, len(filelist))
        return os.path.join(dir, filelist[0])

    def _load_model(self):
        raise NotImplementedError
