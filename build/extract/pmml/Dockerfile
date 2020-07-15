FROM cargo.caicloud.xyz/library/centos:7.4.1708

ARG openscore_url=https://github.com/openscoring/openscoring/releases/download/2.0.1/openscoring-server-executable-2.0.1.jar
ARG ormb_version=v0.0.3
ARG ormb_tar_filename=ormb_0.0.3_Linux_x86_64.tar.gz

RUN rpm  --import  /etc/pki/rpm-gpg/RPM-GPG-KEY-CentOS-7           \
    && yum install -y java-11-openjdk python3 wget                 \
    && pip3 install --no-cache-dir requests loguru pyyaml requests \  
    && wget https://github.com/kleveross/ormb/releases/download/$ormb_version/$ormb_tar_filename \
    && tar -xvf $ormb_tar_filename -C /usr/local/bin   \
    && rm -rf $ormb_tar_filename                       \
    && yum clean all 

COPY application.conf /opt/openscoring/

ADD ${openscore_url} /opt/openscoring
COPY scripts/shell/run.sh /scripts/run.sh
COPY scripts/extract  /scripts

WORKDIR /scripts

ENV EXTRACTOR=pmml