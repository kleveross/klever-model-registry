# Use distroless as minimal base image to package the modeljob-operator binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot

WORKDIR /
COPY bin/klever-modeljob-operator  /klever-modeljob-operator
USER nonroot:nonroot

ENTRYPOINT ["/klever-modeljob-operator"]
