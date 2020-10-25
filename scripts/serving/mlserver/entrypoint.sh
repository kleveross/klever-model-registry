#!/bin/bash

python3 /opt/wrapper/preprocessor.py

if [ "$MODEL_FORMAT" = "MLlib" ];then
    mlservermllib start $MODEL_STORE
else
    mlserver start $MODEL_STORE
fi
