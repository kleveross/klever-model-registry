FROM cargo.caicloud.xyz/library/centos:7.4.1708

ARG url=https://github.com/openscoring/openscoring/releases/download/2.0.1/openscoring-server-executable-2.0.1.jar

RUN rpm  --import  /etc/pki/rpm-gpg/RPM-GPG-KEY-CentOS-7    \
    && yum install -y java-11-openjdk python3 wget          \
    && pip3 install --no-cache-dir requests loguru pyyaml   \  
    && yum clean all 

COPY wrapper /opt/openscoring/wrapper

COPY openscoring/docker/application.conf /opt/openscoring/

ADD ${url} /opt/openscoring

WORKDIR /opt/openscoring

EXPOSE 8000

ENTRYPOINT ["/opt/openscoring/wrapper/openscoring.sh"]
