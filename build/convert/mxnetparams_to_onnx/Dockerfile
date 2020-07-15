FROM python:3.6-slim

ENV LC_ALL="C.UTF-8" \
  LANG="C.UTF-8"

ARG ormb_version=v0.0.3
ARG ormb_tar_filename=ormb_0.0.3_Linux_x86_64.tar.gz

RUN apt update -y \
    && apt install libgomp1 wget -y \
    && pip config set global.index-url https://mirrors.aliyun.com/pypi/simple/ \
    && pip install --no-cache-dir grpcio grpcio-tools onnx==1.3.0 mxnet==1.5.1 pyyaml\
    && wget https://github.com/kleveross/ormb/releases/download/$ormb_version/$ormb_tar_filename \
    && tar -xvf $ormb_tar_filename -C /usr/local/bin \
    && rm -rf $ormb_tar_filename
    
#Set timezone
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

COPY scripts/shell/run.sh /scripts/run.sh
COPY scripts/convert  /scripts/convert
WORKDIR /scripts

