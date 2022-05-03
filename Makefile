VERSION=`git rev-parse --short HEAD`
flags=-ldflags="-s -w -X main.version=${VERSION}"
# flags=-ldflags="-s -w -extldflags -static"

all: build

build:
	GODEBUG=netdns=go CGO_ENABLED=0 go clean; rm -rf pkg; go build ${flags}

build_debug:
	go clean; rm -rf pkg; go build ${flags} -gcflags="-m -m"

build_all: build_osx build_osx_arm64 build_amd64 build)_arm64 build_power8 build_windows

build_osx:
	go clean; rm -rf pkg wmstats_osx; GOOS=darwin go build ${flags}
	mv wmstats wmstats_osx_amd64

build_osx_arm64:
	go clean; rm -rf pkg wmstats_osx; GOOS=darwin go build ${flags}
	mv wmstats wmstats_osx_aarch64

build_amd64:
	go clean; rm -rf pkg wmstats_linux; GOOS=linux go build ${flags}
	mv wmstats wmstats_amd64

build_power8:
	go clean; rm -rf pkg wmstats_power8; GOARCH=ppc64le GOOS=linux go build ${flags}
	mv wmstats wmstats_ppc64le

build_arm64:
	go clean; rm -rf pkg wmstats_arm64; GOARCH=arm64 GOOS=linux go build ${flags}
	mv wmstats wmstats_arm64

build_windows:
	go clean; rm -rf pkg wmstats.exe; GOARCH=amd64 GOOS=windows go build ${flags}
	mv wmstats wmstats_windows

install:
	go install

clean:
	go clean; rm -rf pkg

test : test1

test1:
	cd test; go test
