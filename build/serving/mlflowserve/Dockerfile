FROM continuumio/miniconda3

RUN apt-get update && \
    # install prequired modules to support install of mlflow and related components
    apt-get install -y default-libmysqlclient-dev build-essential curl \
    # cmake and protobuf-compiler required for onnx install
    cmake protobuf-compiler

ARG url=https://github.com/mlflow/mlflow/archive/v1.11.0.tar.gz

ENV MODEL_STORE /mnt/models
ENV PORT 8000

ADD ${url} /
RUN tar -zxvf v1.11.0.tar.gz && \
    mv mlflow-1.11.0 /app

WORKDIR /app

RUN pip install -r dev-requirements.txt --no-cache-dir && \
    pip install -r test-requirements.txt --no-cache-dir && \
    # install mlflow in editable form
    pip install --no-cache-dir -e . && \
    # mkdir required to support install openjdk-11-jre-headless
    mkdir -p /usr/share/man/man1 && apt-get install -y openjdk-11-jre-headless && \
    # install npm for node.js support
    curl -sL https://deb.nodesource.com/setup_10.x | bash - && \
    apt-get update && apt-get install -y nodejs && \
    cd mlflow/server/js && \
    npm install && \
    npm run build


COPY scripts/serving/mlflowserver/entrypoint.sh .

ENTRYPOINT [ "entrypoint.sh" ]
