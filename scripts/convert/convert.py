import os
import argparse

from converter import Converter

if __name__ == "__main__":

    parser = argparse.ArgumentParser(description='Need model path.')
    parser.add_argument(
        '--input_dir',
        help='path to the directory of model needs to be converted',
        dest='input_dir')

    parser.add_argument('--output_dir',
                        help='path to the directory of exported model',
                        dest='output_dir')
    args = parser.parse_args()
    converter = Converter(args.input_dir, args.output_dir)
    converter.convert()
