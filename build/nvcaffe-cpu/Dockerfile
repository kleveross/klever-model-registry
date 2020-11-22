FROM ubuntu:18.04 AS caffe

# Caffe
ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y --no-install-recommends wget curl  build-essential git gfortran libgflags-dev \
      libboost-filesystem-dev libboost-python-dev libboost-system-dev libboost-thread-dev libboost-regex-dev \
      libgoogle-glog-dev libhdf5-serial-dev libleveldb-dev liblmdb-dev libopencv-dev libsnappy-dev python3-wheel python3-setuptools \
      python3-all-dev python3-dev python3-h5py python3-matplotlib python3-numpy python3-opencv python3-pil  \
      python3-pip python3-pydot python3-scipy python3-skimage python3-sklearn libturbojpeg libturbojpeg-dev \
      doxygen libopenblas-dev  libprotobuf-dev protobuf-compiler&& \
    ln -sf /usr/bin/python3 /usr/bin/python &&\
    ln -sf /usr/bin/pip3 /usr/bin/pip &&\
    git clone https://github.com/NVIDIA/caffe.git /usr/src/caffe -b 'v0.16.5' && \
    sed -i '2c numpy>=1.7.1' /usr/src/caffe/python/requirements.txt &&\
    sed -i '51c set(python_version "3" CACHE STRING "Specify which Python version to use")' /usr/src/caffe/CMakeLists.txt && \
    pip install -r /usr/src/caffe/python/requirements.txt 
RUN wget -O cmake.sh https://github.com/Kitware/CMake/releases/download/v3.14.0-rc1/cmake-3.14.0-rc1-Linux-x86_64.sh && \
        bash cmake.sh --skip-license && rm -f cmake.sh
RUN cd /usr/src/caffe && \
    mkdir build && \
    cd build && \
    cmake .. -DCPU_ONLY=1 -DBLAS=open -DCMAKE_INSTALL_PREFIX=/usr/local/caffe && \
    make -j"$(nproc)" && \
    make install

FROM ubuntu:18.04

ENV PYTHONPATH=/usr/local/python:$PYTHONPATH

COPY --from=caffe /usr/local/caffe /usr/local

ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y --no-install-recommends gcc libc6-dev python-dev \
        libboost-python-dev libboost-system-dev libboost-thread-dev libboost-regex-dev libturbojpeg \
        libgoogle-glog-dev libgflags-dev libhdf5-serial-dev libleveldb-dev liblmdb-dev \
        libopencv-dev python-opencv libopenblas-dev libprotobuf-dev protobuf-compiler python3-wheel python3-setuptools \
        ca-certificates \
        curl \
        python3 python3-pip && \
    ln -sf /usr/bin/python3 /usr/bin/python &&\
    ln -sf /usr/bin/pip3 /usr/bin/pip &&\
    rm -rf /var/lib/apt/lists/*


RUN pip install -U --no-cache-dir six matplotlib numpy scipy networkx pillow scikit-image protobuf==3.2.0 pip