#!/bin/bash

mlflow models serve -m $MODEL_STORE -h 0.0.0.0 -p $PORT