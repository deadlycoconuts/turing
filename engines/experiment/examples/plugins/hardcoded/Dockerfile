FROM golang:1.22-alpine as builder

ARG binary_name="plugin"
ARG project_root=github.com/caraml-dev/turing/engines/experiment/examples/plugins/hardcoded/cmd
ARG GOOS="linux"
ARG GOARCH="amd64"

ENV GO111MODULE=on \
    GOOS=${GOOS} \
    GOARCH=${GOARCH}

WORKDIR /app
COPY . .

RUN go build \
    -mod=vendor \
    -o ./bin/${binary_name} \
    -v ${project_root}

FROM alpine:3.17

ARG turing_user="turing"
ARG turing_user_group="app"
ARG binary_name="plugin"

RUN addgroup -S ${turing_user_group} \
    && adduser -S ${turing_user} -G ${turing_user_group} -H

ENV PLUGIN_NAME ""
ENV PLUGINS_DIR "/app/plugins"

COPY --chown=${turing_user}:${turing_user_group} --from=builder /app/bin/${binary_name} /go/bin/plugin

CMD ["sh", "-c", "cp /go/bin/plugin ${PLUGINS_DIR}/${PLUGIN_NAME:?variable must be set}"]