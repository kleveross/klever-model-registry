import os
import argparse
import json
import yaml
import logging

from extractor import Extractor

logging.basicConfig(
    format='[%(levelname).1s%(asctime)s\t%(name)s] %(message)s',
    datefmt='%m%d %I:%M:%S',
    level=logging.INFO)


def update_yaml(dir, res_dict):
    with open(os.path.join(dir, 'ormbfile.yaml'), 'r') as f:
        data = yaml.safe_load(f)

    with open(os.path.join(dir, 'ormbfile.yaml'), 'w') as f:
        layers = res_dict['Operators']

        if 'signature' not in data:
            data['signature'] = {}

        data['signature']['layers'] = dict(layers)

        if len(res_dict['Inputs']) != 0:
            data['signature']['inputs'] = res_dict['Inputs']

        if len(res_dict['Outputs']) != 0:
            data['signature']['outputs'] = res_dict['Outputs']
        logging.info('save ormbfile.yaml: ' + json.dumps(data))
        yaml.dump(data, f)


if __name__ == '__main__':

    parser = argparse.ArgumentParser(description='Process some the path.')
    parser.add_argument('-d', metavar='DIR', help='path to model', dest='dir')

    args = parser.parse_args()
    extractor = Extractor(path=args.dir)

    try:
        res_dict = extractor.extract()
        logging.info('origin data from extractor: ' + json.dumps(res_dict))
        update_yaml(args.dir, res_dict)
    except Exception as e:
        logging.error(e)
