ARG BUILD_IMAGE=golang:1.21-alpine
FROM ${BUILD_IMAGE} as Build

COPY . . 

RUN GOPATH= go build -o /main

FROM alpine:latest

COPY --from=Build /main /main

ENTRYPOINT [ "/main"]
