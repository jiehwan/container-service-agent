package csaapi

import (
	"fmt"
	"github.com/tv42/httpunix"

	"log"
	"net/http"
	"net/http/httputil"
	"time"
)


const (
	ContainerServiceSocket string = "/var/run/container_service.sock"
	DockerLauncherSocket string = "/var/run/docker_launcher.sock"
)


type ContainerInfo struct {
	ContainerID     string `json:"container_id"`
	ContainerStatus string `json:"container_status"`
}
type ContainerLists struct {
	Cmd            string          `json:"cmd"`
	ContainerCount int             `json:"container_count"`
	Container      []ContainerInfo `json:"container"`
}

/*
// go-to-json output is follows.. but there is problem during init.
type ContainerLists struct {
	Cmd string `json:"cmd"`
	ContainerCount int `json:"container_count"`
	Container []struct {
		ContainerID string `json:"container_id"`
		ContainerStatus string `json:"container_status"`
	} `json:"container"`
}
*/
type ContainerService interface {
	GetContainersInfo() (ContainerLists, error)
}

/*
	STUB
*/
func GetContainersInfo() (ContainerLists, error) {

	send := ContainerLists{
		Cmd:            "getContainerLists",
		ContainerCount: 2,
		Container: []ContainerInfo{
			{
				ContainerID:     "api-1111",
				ContainerStatus: "running",
			},
			{
				ContainerID:     "api-2222",
				ContainerStatus: "exited",
			},
		},
	}

	u := &httpunix.Transport{
		DialTimeout:           100 * time.Millisecond,
		RequestTimeout:        1 * time.Second,
		ResponseHeaderTimeout: 1 * time.Second,
	}
	u.RegisterLocation("myservice", "/var/run/container_service.sock")

	var client = http.Client{
		Transport: u,
	}

	resp, err := client.Get("http+unix://myservice/getContainersInfo")
	if err != nil {
		log.Fatal(err)
	}
	buf, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", buf)
	resp.Body.Close()

	return send, nil
}
