#!/bin/bash

python3 /opt/wrapper/preprocessor.py
mlserver start $MODEL_STORE
