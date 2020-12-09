FROM ghcr.io/kleveross/nvcaffe:cpu-0.16.5

ENV LC_ALL="C.UTF-8" \
  LANG="C.UTF-8"

ARG ORMB_VERSION=0.0.8
ARG ORMB_TAG=v${ORMB_VERSION}
ARG ORMB_TAR_FILENAME=ormb_${ORMB_VERSION}_Linux_x86_64.tar.gz

RUN apt update && apt install -y wget && \
    pip install --no-cache-dir \
                grpcio==1.26.0 \
                grpcio-tools==1.26.0 \ 
                pyyaml && \
    wget https://github.com/caicloud/ormb/releases/download/$ORMB_TAG/$ORMB_TAR_FILENAME && \
    tar -xvf $ORMB_TAR_FILENAME -C /usr/local/bin && \
    rm -rf $ORMB_TAR_FILENAME

#Set timezone
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

COPY scripts/shell/*.sh /scripts/
COPY scripts/extract  /scripts

ENV EXTRACTOR=caffemodel
ENV SOURCE_FORMAT=CaffeModel
ENV FORMAT=CaffeModel

ENTRYPOINT ["sh","-c","/scripts/run.sh"]
