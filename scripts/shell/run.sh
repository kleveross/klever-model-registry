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
echo "ORMB domain: $SERVER_ORMB_DOMAIN"
echo "ORMB username: $SERVER_ORMB_USERNAME"
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
ormb login  --insecure $SERVER_ORMB_DOMAIN -u $SERVER_ORMB_USERNAME -p $SERVER_ORMB_PASSWORD
checkOrExit $? $ormb_login_err

if [ $dst_tag == "empty" ]
then

    case $format in
    PMML )
        java -Dconfig.file=/opt/openscoring/application.conf -jar /opt/openscoring/openscoring-server-executable-2.0.1.jar --port 8080 &
        sleep 15s;
        ;;
    esac

    # execute python script to extract
    python3 /scripts/extract.py -d $input_dir
    checkOrExit $? $ormb_run_task_err

else
    python3 /scripts/convert.py --input_dir=$input_dir --output_dir=$output_dir
    checkOrExit $? $ormb_run_task_err
    
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
