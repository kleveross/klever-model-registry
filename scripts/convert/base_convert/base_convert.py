import os
import yaml

class BaseConvert(object):
    def __init__(self):
        pass
    
    def read_ormbfile(self, dir):
        data = {}
        with open(os.path.join(dir, "ormbfile.yaml"), 'r') as f:
            print("ormbfile path: ", dir)
            data = yaml.load(f)
            print("ormtfile content: ", data)

        return data

    def write_output_ormbfile(self, input_dir, output_dir):
        data = self.read_ormbfile(input_dir)
        output_ormbfile = {}
        if "author" in data:
            output_ormbfile["author"] = data["author"]

        output_ormbfile["format"] = os.environ["FORMAT"]
        
        with open(os.path.join(output_dir, "ormbfile.yaml"), "w+") as f:
            yaml.dump(output_ormbfile, f)


