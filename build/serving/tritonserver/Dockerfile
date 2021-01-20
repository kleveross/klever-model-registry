FROM nvcr.io/nvidia/tritonserver:20.09-py3

COPY scripts/serving/tensorrt /

RUN apt update \
    && apt install -y python3.6 \
    && apt install -y python3-pip \
    && pip3 install -r /requirements.txt \
    && rm /requirements.txt

COPY scripts/serving/wrapper /opt/wrapper

ENTRYPOINT ["/entrypoint.sh"]
