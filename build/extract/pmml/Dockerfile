FROM centos:7.4.1708

ARG openscore_url=https://github.com/openscoring/openscoring/releases/download/2.0.1/openscoring-server-executable-2.0.1.jar
ARG ORMB_VERSION=0.0.8
ARG ORMB_TAG=v${ORMB_VERSION}
ARG ORMB_TAR_FILENAME=ormb_${ORMB_VERSION}_Linux_x86_64.tar.gz

RUN rpm  --import  /etc/pki/rpm-gpg/RPM-GPG-KEY-CentOS-7           && \
    yum install -y java-11-openjdk python3 wget                 && \
    pip3 install --no-cache-dir requests loguru pyyaml requests && \  
    wget https://github.com/caicloud/ormb/releases/download/$ORMB_TAG/$ORMB_TAR_FILENAME && \
    tar -xvf $ORMB_TAR_FILENAME -C /usr/local/bin && \
    rm -rf $ORMB_TAR_FILENAME && \
    yum clean all 

COPY build/extract/pmml/application.conf /opt/openscoring/

ADD ${openscore_url} /opt/openscoring
COPY scripts/shell/ /scripts/
COPY scripts/extract/  /scripts/

ENV EXTRACTOR=pmml
ENV SOURCE_FORMAT=PMML
ENV FORMAT=PMML

ENTRYPOINT ["sh","-c","/scripts/run.sh"]