#!/bin/bash

run_task_err=10003

source_format=$SOURCE_FORMAT
format=$FORMAT


function checkOrExit() {
    if [ $1 != 0 ];then
        echo "exit code: $2"
        exit $2
    fi
}

# mkdir for output
output_dir=/tmp/test_model
mkdir -p ${output_dir}/model

if [ $source_format == $format ]
then

    case $format in
    PMML )
        java -Dconfig.file=/opt/openscoring/application.conf -jar /opt/openscoring/openscoring-server-executable-2.0.1.jar --port 8080 &
        sleep 15s;
        ;;
    esac
    checkOrExit $? $run_task_err

    # execute python script to extract
    case $format in 
    CaffeModel )
        EXTRACTOR=caffemodel python /scripts/extract.py -d /models/Caffe-model/
        ;;
    GraphDef )
        EXTRACTOR=graphdef python /scripts/extract.py -d /models/GraphDef-model/
        ;;
    H5 )
        EXTRACTOR=h5 python /scripts/extract.py -d /models/H5-model/
        ;;
    MXNetParams )
        EXTRACTOR=mxnetparams python /scripts/extract.py -d /models/MXNetParams-model/
        ;;  
    NetDef )
        EXTRACTOR=netdef python /scripts/extract.py -d /models/NetDef-model/
        ;; 
    ONNX )
        EXTRACTOR=onnx python /scripts/extract.py -d /models/ONNX-model/
        ;;  
    SavedModel )
        EXTRACTOR=savedmodel python /scripts/extract.py -d /models/SavedModel-model/
        ;;    
    TorchScript )
        EXTRACTOR=torchscript python /scripts/extract.py -d /models/TorchScript-model/
        ;;  
    esac
    checkOrExit $? $run_task_err

else
    case $format in 
    ONNX )
        case $source_format in 
        MXNetParams )
            INPUTS='[{"name":"data","size":[1,3,224,224],"dType":"float32"}]'  python /scripts/convert.py --input_dir=/models/MXNetParams-model/ --output_dir=$output_dir
            ;;
        NetDef )
            INPUTS='[{"name":"data","size":[1,3,224,224],"dType":"float32"}]'  python /scripts/convert.py --input_dir=/models/NetDef-model/ --output_dir=$output_dir
            ;;
        esac
        ;;
    NetDef )
        USING_ORMBFILE=True  python /scripts/convert.py --input_dir=/models/Caffe-model/ --output_dir=$output_dir
        ;;
    SavedModel )
        USING_ORMBFILE=True  python /scripts/convert.py --input_dir=/models/H5-model/ --output_dir=$output_dir
        ;;
    esac
    checkOrExit $? $run_task_err
fi