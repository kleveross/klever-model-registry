import os

from caffe.proto import caffe_pb2
from google.protobuf import text_format

from caffe2.python.caffe_translator import TranslateModel
from caffe2.proto import caffe2_pb2
from caffe2.python import core, utils, workspace

from .base_convert import BaseConverter

INIT_NET = 'init_net.pb'
PREDICT_NET = 'predict_net.pb'

MODEL_EXTENSION = '.prototxt'
WEIGHTS_EXTENSION = '.caffemodel'


def ConvertTensorProtosToInitNet(net_params, input_names):
    # rewrite to support multi inputs
    """Takes the net_params returned from TranslateModel, and wrap it as an
    init net that contain GivenTensorFill.
    This is a very simple feature that only works with float tensors, and is
    only intended to be used in an environment where you want a single
    initialization file - for more complex cases, use a db to store the
    parameters.
    """
    init_net = caffe2_pb2.NetDef()
    for tensor in net_params.protos:
        if len(tensor.float_data) == 0:
            raise RuntimeError(
                "Only float tensors are supported in this util.")
        op = core.CreateOperator(
            "GivenTensorFill", [], [tensor.name],
            arg=[
                utils.MakeArgument("shape", list(tensor.dims)),
                utils.MakeArgument("values", tensor.float_data)
            ])
        init_net.op.extend([op])
    for input_name in input_names:
        init_net.op.extend(
            [core.CreateOperator("ConstantFill", [], [input_name], shape=[1])])
    return init_net


class CaffeToCaffe2(BaseConverter):
    def _load_model(self):
        self.caffemodel = self._find_with_extension(WEIGHTS_EXTENSION)
        self.prototext = self._find_with_extension(MODEL_EXTENSION)

    def _convert(self):

        caffenet = caffe_pb2.NetParameter()
        caffenet_pretrained = caffe_pb2.NetParameter()

        input_proto = self.prototext
        input_caffemodel = self.caffemodel

        ### Output
        output_init_net = os.path.join(self.output_dir, 'model', INIT_NET)
        output_predict_net = os.path.join(self.output_dir, 'model',
                                          PREDICT_NET)

        with open(input_proto) as f:
            text_format.Merge(f.read(), caffenet)
        with open(input_caffemodel, 'rb') as f:
            caffenet_pretrained.ParseFromString(f.read())

        net, pretrained_params = TranslateModel(caffenet,
                                                caffenet_pretrained,
                                                is_test=True,
                                                input_dims=[])

        # multi inputs and output
        all_inputs = [op_input for op in net.op for op_input in op.input]
        all_outputs = [op_out for op in net.op for op_out in op.output]

        all_inputs = set(all_inputs)
        all_outputs = set(all_outputs)

        param_input = [param.name for param in pretrained_params.protos]
        param_set = set(param_input)

        real_inputs = [
            op_input for op_input in all_inputs
            if op_input not in all_outputs and op_input not in param_set
        ]
        real_outputs = [
            op_out for op_out in all_outputs if op_out not in all_inputs
        ]

        net.external_input.extend(real_inputs)
        net.external_input.extend(param_input)
        net.external_output.extend(real_outputs)

        init_net = ConvertTensorProtosToInitNet(pretrained_params, real_inputs)

        with open(output_predict_net, 'wb') as f:
            f.write(net.SerializeToString())
        with open(output_predict_net + 'txt', 'w') as f:
            f.write(str(net))
        with open(output_init_net, 'wb') as f:
            f.write(init_net.SerializeToString())


if __name__ == '__main__':
    convert = CaffeToCaffe2()
    convert.convert()
