import requests
import os
from .base_extract import BaseExtrctor
import json

class PMMLExtractor(BaseExtrctor):
    def _extrat_params(self):
        for file in os.listdir(os.path.join(self.dir, "model")):
            if file.endswith(".pmml"):
                file_path = os.path.join(self.dir, "model", file)
                with open(file_path) as f:
                    content = f.read()
                    resp = requests.put("http://127.0.0.1:8080/openscoring/model/extract", data=content)
                    return json.loads(resp.text.encode("utf-8").decode("utf-8"))
                break
        
    def _extract_inputs(self):
        data = self._extrat_params()
        return data["schema"]["inputFields"]

    def _extract_outputs(self):
        data = self._extrat_params()
        return data["schema"]["outputFields"]
    
    def _extract_ops(self):
        return []

    def _load_model(self):
        pass