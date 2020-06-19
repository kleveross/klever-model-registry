# Use distroless as minimal base image to package the model-registry-controller binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot

WORKDIR /
COPY bin/model-registry-controller  /model-registry-controller
USER nonroot:nonroot

ENTRYPOINT ["/model-registry-controller"]
