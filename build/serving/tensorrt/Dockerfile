FROM nvcr.io/nvidia/tensorrtserver:19.10-py3

COPY scripts/serving/tensorrt /

RUN pip3 install -r /requirements.txt && rm /requirements.txt

COPY scripts/serving/wrapper /opt/wrapper

ENTRYPOINT ["/entrypoint.sh"]
