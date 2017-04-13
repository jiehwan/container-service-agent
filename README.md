## container-service-agent

container-service-agent is designed to manage between Tizen platform and server.

### Developer Quick-Start

To build the daemon , the following build system dependencies are required:

* go 1.7.5 or above
* gb tool due to library dependecy

#### go 1.7.5

```
$ wget https://storage.googleapis.com/golang/go1.7.5.linux-amd64.tar.gz

```
If you extract the file and see the 'go' folder.
Copy 'go' folder into '/usr/local/go'
Set up the GOROOT, GOPATH, PATH

```
$ export PATH=$PATH:/usr/local/go/bin/
$ export GOPATH=$(go env GOPATH)
$ export PATH=$PATH:$(go env GOPATH)/bin
```

#### gb

```
$ go get github.com/constabulary/gb/ ...
$ gb info
```

#### build

```
$ mkae build
```
/bin/ folder is created and **two binaries** you can see in the folder.
**container-service** is client which can receive a command from server
**container-serviced** is main daemon to check request form container-service and request to docker-launcher

