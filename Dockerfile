# stage -- build binary
FROM golang:1.19.9-alpine3.17 AS stage
WORKDIR /workspace
COPY . .
RUN go env -w GOPROXY="https://goproxy.cn,direct" && go build -o garden main.go

# build iamge
FROM alpine:3.17
WORKDIR /workspace
COPY --from=stage /workspace/garden .
COPY --from=stage /workspace/db/migration ./db/migration
COPY --from=stage /workspace/start.sh .
COPY --from=stage /workspace/wait-for.sh .

# TODO set port by env
ENTRYPOINT ["/workspace/wait-for.sh", "pgsql:5432", "--", "/workspace/start.sh"]
CMD ["/workspace/garden", "server"]

