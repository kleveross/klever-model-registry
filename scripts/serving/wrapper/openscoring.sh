#!/bin/bash

python3 /opt/openscoring/wrapper/preprocessor.py

if [ -f "${MODEL_STORE}/../${SERVING_NAME}/1/model.pmml" ];then
    while true;do
        result=$(curl -X PUT --data-binary @${MODEL_STORE}/../${SERVING_NAME}/1/model.pmml -H "Content-type: text/xml" http://localhost:8000/openscoring/model/${SERVING_NAME} -o /dev/null -s -w %{http_code})
        if [ $result -eq 201 ]\
            || [ $result -eq 200 ];then
            break
        else
            sleep 2
        fi
    done &

    java -Dconfig.file=/opt/openscoring/application.conf -jar /opt/openscoring/openscoring-server-executable-2.0.1.jar --port 8000
else
    exit 1
fi
