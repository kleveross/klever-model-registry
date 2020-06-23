import os
import argparse
import yaml

from extractor import Extractor

ID = os.environ.get('ANALYZE_ID', 'NULL')
IP = os.environ.get('ANALYZE_HOST', '0.0.0.0:8081')


def dict2pb(d, pb):

    for Input in d['Inputs']:
        temp_input = pb.Inputs.add()
        temp_input.Name = Input.get('signatureConst', Input['name'])
        temp_input.DataType = Input['dataType']
        for dim in Input['dims']:
            temp_input.Dims.append(dim)

    for Output in d['Outputs']:
        temp_output = pb.Outputs.add()
        temp_output.Name = Output.get('signatureConst', Output['name'])
        temp_output.DataType = Output['dataType']
        for dim in Output['dims']:
            temp_output.Dims.append(dim)

    for name, num in d['Operators'].items():
        pb.Operators[name] = num


def update_yaml(dir, res_dict):
    print(res_dict)

    data = dict
    with open(os.path.join(dir, "ormbfile.yaml"), 'r') as f:
        data = yaml.load(f)

    with open(os.path.join(dir, "ormbfile.yaml"), 'w') as f:
        data["framework"] = os.environ["FRAMEWORK"]
        data["format"] = os.environ["FORMAT"]
        layersMap = dict(res_dict["Operators"])
        layers = []
        for (k, v) in  layersMap.items(): 
            layers.append({"name": k})
        data["signature"] = {
            "layers": layers
        }

        if len(res_dict["Inputs"]) != 0:
            data["signature"]["inputs"] = res_dict["Inputs"]

        if len(res_dict["Outputs"]) != 0:
            data["signature"]["outputs"] = res_dict["Outputs"]
        
        yaml.dump(data, f)

if __name__ == "__main__":

    parser = argparse.ArgumentParser(description='Process some the path.')
    parser.add_argument('-d', metavar='DIR', help='path to model', dest='dir')

    args = parser.parse_args()
    extractor = Extractor(path=args.dir)

    try:
        res_dict = extractor.extract()
        print(res_dict)
        update_yaml(args.dir, res_dict)
    except Exception as e:
        print(str(e))
