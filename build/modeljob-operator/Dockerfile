# Use distroless as minimal base image to package the modeljob-operator binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot

WORKDIR /
COPY bin/modeljob-operator  /modeljob-operator
COPY /crds /crds
USER nonroot:nonroot

ENTRYPOINT ["/modeljob-operator"]
