FROM python:3.6-slim

ENV LC_ALL="C.UTF-8" \
  LANG="C.UTF-8"

ARG ORMB_VERSION=0.0.8
ARG ORMB_TAG=v${ORMB_VERSION}
ARG ORMB_TAR_FILENAME=ormb_${ORMB_VERSION}_Linux_x86_64.tar.gz
  
RUN apt update && apt install -y wget libgomp1 git && \
    pip install --no-cache-dir \
                tools \
                grpcio-tools==1.26.0 \
                https://download.pytorch.org/whl/cpu/torch-1.7.0%2Bcpu-cp36-cp36m-linux_x86_64.whl \ 
                future \ 
                pyyaml && \
    wget https://github.com/caicloud/ormb/releases/download/$ORMB_TAG/$ORMB_TAR_FILENAME && \
    tar -xvf $ORMB_TAR_FILENAME -C /usr/local/bin && \
    rm -rf $ORMB_TAR_FILENAME

#Set timezone
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

COPY scripts/shell/*.sh /scripts/
COPY scripts/convert  /scripts

RUN git clone https://github.com/NVIDIA/caffe.git /usr/src/caffe  -b 'v0.17.3'  &&  \
    mkdir -p /scripts/caffe/proto &&  touch /scripts/caffe/__init__.py && \ 
    touch /scripts/caffe/proto/__init__.py && \
    python -m grpc_tools.protoc -I/usr/src/caffe/src/caffe/proto \
                              --python_out=/scripts/caffe/proto \ 
                              --grpc_python_out=/scripts/caffe/proto /usr/src/caffe/src/caffe/proto/*.proto

ENV SOURCE_FORMAT=CaffeModel
ENV FORMAT=NetDef

ENTRYPOINT ["sh","-c","/scripts/run.sh"]

