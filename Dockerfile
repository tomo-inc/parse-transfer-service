FROM  golang:1.23.8-alpine3.21 AS builder
LABEL stage=builder
RUN apk add --no-cache git
RUN mkdir /go/src/app
WORKDIR /go/src/app
COPY ./ ./
RUN export Branch=$(git symbolic-ref --short -q HEAD) && \
    export Commit=$(git rev-parse --short HEAD) && \
    export Date=$(date "+%Y-%m-%d %H:%M:%S") && \
    export Author=$(git log --pretty=format:%an $(Commit) -1) && \
    export Email=$(git log --pretty=format:%ae $(Commit) -1) && \
    export GoVersion=$(go version) && \
    CGO_ENABLED=0 go build -o parse -a -installsuffix cgo \
    -ldflags "-X 'github.com/tomo.inc/parse-transfer-service/cmd/version.branch=${Branch}' \
        -X 'github.com/tomo.inc/parse-transfer-service/cmd/version.author=${Author}' \
        -X 'github.com/tomo.inc/parse-transfer-service/cmd/version.date=${Date}' \
        -X 'github.com/tomo.inc/parse-transfer-service/cmd/version.commit=${Commit}' \
        -X 'github.com/tomo.inc/parse-transfer-service/cmd/version.goVersion=${GoVersion}'"

FROM alpine:3.21
WORKDIR /root/
COPY --from=builder /go/src/app/parse .
ENTRYPOINT ["/root/parse"]
