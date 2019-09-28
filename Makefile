BINARY_PATH=bin
GPM_BINARY=gpm
GPM_PATH=./

build: $(GPM_BINARY)

$(GPM_BINARY):
	go build -o ${BINARY_PATH}/${GPM_BINARY} ${GPM_PATH}

clean:
	rm -rf bin && go clean