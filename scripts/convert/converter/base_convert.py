import os
import json
import yaml

INPUTS_ENV = "INPUTS"
OUTPUTS_ENV = "OUTPUTS"
USING_ORMBFILE_ENV = "USING_ORMBFILE"


class BaseConverter(object):
    def __init__(self, input_dir, output_dir):
        self.input_dir = input_dir
        self.output_dir = output_dir
        self.using_ormbfile = USING_ORMBFILE_ENV in os.environ
        self._parase_modelfile()

    def _parase_modelfile(self):
        with open(os.path.join(self.input_dir, 'ormbfile.yaml'), 'r') as f:
            data = yaml.safe_load(f)

        self.author = data.get('author', None)

        # maybe can read inputs in env
        if self.using_ormbfile:
            self.input_value = data['signature'].get('inputs', [])
            self.output_value = data['signature'].get('outputs', [])
        else:
            self.input_value = json.loads(os.getenv(INPUTS_ENV, '[]'))
            self.output_value = json.loads(os.getenv(OUTPUTS_ENV, '[]'))

        assert len(self.input_value) != 0, 'expected models inputs'

    def _write_output_ormbfile(self):
        output_ormbfile = dict()
        if self.author:
            output_ormbfile['author'] = self.author

        output_ormbfile['format'] = os.environ['FORMAT']

        with open(os.path.join(self.output_dir, 'ormbfile.yaml'), 'w') as f:
            yaml.safe_dump(output_ormbfile, f)

    def _find_with_extension(self, extension):
        dir = os.path.join(self.input_dir, 'model')
        filelist = list(
            filter(lambda f: f.endswith(extension) and not f.startswith('.'),
                   os.listdir(dir)))
        assert (len(filelist) == 1), 'expected one %s file,but found %s' % (
            extension, len(filelist))
        return os.path.join(dir, filelist[0])

    def _find_with_name(self, name):
        dir = os.path.join(self.input_dir, 'model')
        filelist = list(filter(lambda f: f == name, os.listdir(dir)))
        assert (len(filelist) == 1
                ), 'expected one %s file,but found %s' % (name, len(filelist))
        return os.path.join(dir, filelist[0])

    def _convert(self):
        raise NotImplementedError

    def _load_model(self):
        raise NotImplementedError

    def convert(self):
        self._load_model()
        self._convert()
        self._write_output_ormbfile()
