FROM python:3.7-slim

ENV MODEL_STORE /mnt/models

RUN apt-get update && \
    apt-get -y --no-install-recommends install \
      libgomp1

# Use MLServer for serving, see https://github.com/SeldonIO/MLServer
WORKDIR /workspace
COPY scripts/serving/mlserver .
RUN pip install .[all]

COPY scripts/serving/wrapper /opt/wrapper
RUN pip install -r /opt/wrapper/requirements.txt && rm /opt/wrapper/requirements.txt

ENTRYPOINT ["/workspace/entrypoint.sh"]
