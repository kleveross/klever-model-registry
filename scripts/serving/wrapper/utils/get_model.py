# -*- coding: utf-8 -*-
import os
import zipfile
import shutil
import yaml
from loguru import logger


@logger.catch()
def check_model(dirPath):
    yaml_file_path = os.path.join(dirPath, 'ormbfile.yaml')
    with open(yaml_file_path, encoding="utf-8") as mf:
        buffer = mf.read()
        manifest = yaml.load(buffer)
    return buffer, manifest

@logger.catch()
def find_xx_file(search_dir, ext, show_all=False):
    all_results = filter(lambda name: name.endswith(ext),
                         os.listdir(search_dir))
    if show_all:
        return list(all_results)
    else:
        return list(all_results)[0]
