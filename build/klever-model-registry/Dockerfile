# Use distroless as minimal base image to package the model-registry binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot

WORKDIR /
COPY bin/klever-model-registry  /klever-model-registry
USER nonroot:nonroot

ENTRYPOINT ["/klever-model-registry"]
