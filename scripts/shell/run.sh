#!/bin/bash

ormb_login_err=10000
ormb_pull_model_err=10001
ormb_export_model_err=10002
ormb_run_task_err=10003
ormb_save_model_err=10004
ormb_push_model_err=10005

src_tag=$SOURCE_MODEL_TAG
dst_tag=$DESTINATION_MODEL_TAG
input_dir=$SOURCE_MODEL_PATH
output_dir=$DESTINATION_MODEL_PATH
source_format=$SOURCE_FORMAT
format=$FORMAT

echo "#####################################################"
echo "model source tag: $src_tag"
echo "model destination tag: $dst_tag"
echo "model input dir: $input_dir"
echo "model output dir: $output_dir"
echo "model format: $format"
echo "ORMB domain: $ORMB_DOMAIN"
echo "ORMB username: $ORMB_USERNAME"
echo "ORMB password: $ORMB_PASSWORD"
echo "#####################################################"


function checkOrExit() {
    if [ $1 != 0 ];then
        echo "exit code: $2"
        exit $2
    fi
}

mkdir -p $input_dir/model
mkdir -p $output_dir/model

# login to harbor.
ormb login  --insecure $ORMB_DOMAIN -u $ORMB_USERNAME -p $ORMB_PASSWORD
checkOrExit $? $ormb_login_err

# pull source model.
ormb pull $src_tag --plain-http
checkOrExit $? $ormb_pull_model_err

# export source model.
ormb export -d $input_dir $src_tag
checkOrExit $? $ormb_export_model_err

if [ $dst_tag == "empty" ]
then

    case $format in
    PMML )
        java -Dconfig.file=/opt/openscoring/application.conf -jar /opt/openscoring/openscoring-server-executable-2.0.1.jar --port 8080 &
        sleep 15s;
        ;;
    NetDef )
        cp $input_dir/model/predict_net.pb $input_dir/model/model.netdef
        cp $input_dir/model/init_net.pb $input_dir/model/init_model.netdef
        ;;
    GraphDef )
        cp $input_dir/model/*.pb  $input_dir/model/model.graphdef
        ;;
    TorchScript )
        cp $input_dir/model/*.pt  $input_dir/model/model.pt
        ;;
    esac

    # execute python script to extract
    python3 /scripts/extract.py -d $input_dir
    checkOrExit $? $ormb_run_task_err

    case $format in
    NetDef )
        rm -rf $input_dir/model/model.netdef
        rm -rf $input_dir/model/init_model.netdef
        ;;
    GraphDef )
        rm -rf $input_dir/model/model.graphdef
        ;;
    TorchScript )
        # rm -rf $input_dir/model/model.pt
        ;;
    esac
else
    case $format in
    ONNX )
        python3 /scripts/convert/convert_mxnet.py --input_dir=$input_dir --output_dir=$output_dir
        ;;
    SavedModel )
        python3 /scripts/convert/convert_keras.py --input_dir=$input_dir --output_dir=$output_dir
        ;;
    NetDef )
        python3 /scripts/convert/convert_caffe.py --input_dir=$input_dir --output_dir=$output_dir # --input_value="[{\"Name\":\"data\",\"Dims\":[1,3,224,224],\"DataType\":\"float32\"}]"
        ;;
    esac
    
fi


if [ $dst_tag == "empty" ]
then
    output_dir=$input_dir
    dst_tag=$src_tag
fi

# save model 
ormb save $output_dir $dst_tag
checkOrExit $? $ormb_save_model_err

# push model to registry
ormb push $dst_tag --plain-http
checkOrExit $? $ormb_push_model_err
