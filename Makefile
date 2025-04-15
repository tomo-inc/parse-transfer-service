Branch=$(shell git symbolic-ref --short -q HEAD)
Commit=$(shell git rev-parse --short HEAD)
#Date=$(shell git log --pretty=format:%cd $(Commit) -1)
Date=$(shell date "+%Y-%m-%d %H:%M:%S")
Author=$(shell git log --pretty=format:%an $(Commit) -1)
shortDate=$(shell git log -1 --format="%at" | xargs -I{} date -d @{} +%Y%m%d)
Email=$(shell git log --pretty=format:%ae $(Commit) -1)
Ver=$(shell echo $(Branch)-$(Commit)-$(shortDate))
GoVersion=$(shell go version )

start: build

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test ./...

build: fmt vet
	go build -ldflags=" \
	-X 'github.com/tomo.inc/parse-transfer-service/cmd/version.branch=$(Branch)' \
	-X 'github.com/tomo.inc/parse-transfer-service/cmd/version.author=$(Author)' \
	-X 'github.com/tomo.inc/parse-transfer-service/cmd/version.date=$(Date)' \
	-X 'github.com/tomo.inc/parse-transfer-service/cmd/version.commit=$(Commit)' \
	-X 'github.com/tomo.inc/parse-transfer-service/cmd/version.goVersion=$(GoVersion)'" \
	-o build/parse
