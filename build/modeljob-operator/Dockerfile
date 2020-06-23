# Use distroless as minimal base image to package the model-registry-controller binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot

WORKDIR /
COPY bin/modeljob-operator  /modeljob-operator
USER nonroot:nonroot

ENTRYPOINT ["/modeljob-operator"]
