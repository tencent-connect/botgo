#执行指令 /usr/bin/make -f Makefile -C ./ build
#打压缩包 zip scf_code.zip * -r
BINARY_NAME=qqbot-svr

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=warn"  -o $(BINARY_NAME) -v

