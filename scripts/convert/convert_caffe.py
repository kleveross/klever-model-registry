import os
import argparse

from caffe.proto import caffe_pb2
from google.protobuf import text_format

from caffe2.python.caffe_translator import TranslateModel
from caffe2.proto import caffe2_pb2
from caffe2.python import core, utils, workspace

from base_convert.base_convert import BaseConvert

INIT_NET = 'init_net.pb'
PREDICT_NET = 'predict_net.pb'


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


if __name__ == '__main__':
    parser = argparse.ArgumentParser(
        description=
        "Utilitity to convert pretrained caffe models to Caffe2 models.")

    parser.add_argument('--input_dir',
                        help='path to the directory of model needs to be converted',
                        dest='input_dir')

    parser.add_argument('--output_dir',
                        help='path to the directory of exported caffe2 model',
                        dest='output_dir')

    parser.add_argument("--remove_legacy_pad",
                        help="Remove legacy pad \
                        (Only works for nets with one input blob)",
                        action="store_true",
                        default=False)
    parser.add_argument("--input_dims",
                        help="Dimension of input blob",
                        nargs='+',
                        type=int,
                        default=[])
    parser.add_argument('--input_value',
                        help='the inputs struct of the model',
                        dest='input_value')
    args = parser.parse_args()

    caffenet = caffe_pb2.NetParameter()
    caffenet_pretrained = caffe_pb2.NetParameter()

    input_proto = ""
    input_caffemodel = ""
    input_model_dir = os.path.join(args.input_dir, "model")
    for file in os.listdir(input_model_dir):
        if file.endswith(".prototxt"):
            input_proto = os.path.join(os.path.join(input_model_dir, file))
        if file.endswith(".caffemodel"):    
            input_caffemodel = os.path.join(os.path.join(input_model_dir, file))

    os.makedirs(args.output_dir, exist_ok=True)
    output_model_dir = os.path.join(args.output_dir, "model")
    output_init_net = os.path.join(output_model_dir, INIT_NET)
    output_predict_net = os.path.join(output_model_dir, PREDICT_NET)

    with open(input_proto) as f:
        text_format.Merge(f.read(), caffenet)
    with open(input_caffemodel, 'rb') as f:
        caffenet_pretrained.ParseFromString(f.read())

    net, pretrained_params = TranslateModel(
        caffenet,
        caffenet_pretrained,
        is_test=True,
        remove_legacy_pad=args.remove_legacy_pad,
        input_dims=args.input_dims)

    ##########################################
    #     ori scripts
    #     Assume there is one input and one output
    #     external_input = net.op[0].input[0]
    #     external_output = net.op[-1].output[0]

    #     net.external_input.extend([external_input])
    #     net.external_input.extend([param.name for param in pretrained_params.protos])
    #     net.external_output.extend([external_output])
    #     init_net = ConvertTensorProtosToInitNet(pretrained_params, external_input)
    ##########################################

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
    
    convert = BaseConvert()
    convert.write_output_ormbfile(args.input_dir, args.output_dir)
