export GOPATH := $(shell gb env | grep GB_SRC_PATH | sed -r "s/\/(\w+)\/src/\/\1/g" | sed -r "s:GB_SRC_PATH=|\"::g")
build:
	gb build all 
clean:
	rm -rf bin && rm -rf pkg
