# Use distroless as minimal base image to package the model-registry binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM debian:stretch

WORKDIR /
COPY bin/model-registry  /model-registry
COPY /crds /crds

ENTRYPOINT ["/model-registry"]
