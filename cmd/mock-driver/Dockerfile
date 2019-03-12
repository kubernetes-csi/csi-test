FROM gcr.io/distroless/static:latest
LABEL maintainers="Kubernetes Authors"
LABEL description="CSI Mock Driver"

# For historic reasons the binary is called "mock" inside the container.
# It's kept that way because some .yaml file might use that name instead
# of relying on the entry point.
COPY ./bin/mock-driver mock
ENTRYPOINT ["/mock"]
