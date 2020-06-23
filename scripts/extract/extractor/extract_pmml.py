import requests
import os
from .base_extract import BaseExtrctor
import json

class PMMLExtractor(BaseExtrctor):
    def _extrat_params(self, name):
        for file in os.listdir(os.path.join(self.dir, "model")):
            if file.endswith(".pmml"):
                file_path = os.path.join(self.dir, "model", file)
                with open(file_path) as f:
                    content = f.read()
                    resp = requests.put("http://127.0.0.1:8080/openscoring/model/extract-"+ name, data=content)
                    data = json.loads(resp.text.encode("utf-8").decode("utf-8"))
                    print(data)
                    return data
                break
        
    def _extract_inputs(self):
        data = self._extrat_params("input")
        inputs = []
        if "schema" not in data:
            return inputs

        if "inputFields" in data["schema"]:
            inputFields = data["schema"]["inputFields"]
            for item in inputFields:
                item["name"] = item["id"]
                item["dtype"] = item["dataType"]
                item["optype"] = item["opType"]
                del item["id"]
                del item["dataType"]
                del item["opType"]
                inputs.append(item)
        return inputs

    def _extract_outputs(self):
        data = self._extrat_params("output")
        outputs = []

        if "schema" in data and "outputFields" in data["schema"]:
            outputFields = data["schema"]["outputFields"]
        
            for item in outputFields:
                item["name"] = item["id"]
                item["dtype"] = item["dataType"]
                item["optype"] = item["opType"]
                del item["id"]
                del item["dataType"]
                del item["opType"]
                outputs.append(item)

        if "schema" in data and "targetFields" in data["schema"]:
            targetFields = data["schema"]["targetFields"]

            for item in targetFields:
                item["name"] = item["id"]
                item["dtype"] = item["dataType"]
                item["optype"] = item["opType"]
                del item["id"]
                del item["dataType"]
                del item["opType"]
                outputs.append(item)

        return outputs
    
    def _extract_ops(self):
        return {}

    def _load_model(self):
        pass