FROM nvcr.io/nvidia/tensorrtserver:19.10-py3

COPY tensorrt/docker/requirements.txt /

RUN pip3 install -r /requirements.txt && rm /requirements.txt

COPY wrapper /opt/wrapper

ENTRYPOINT ["/opt/wrapper/entrypoint.sh"]
