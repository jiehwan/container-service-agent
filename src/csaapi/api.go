package csaapi

import (
	"fmt"
	"encoding/json"
	"github.com/tv42/httpunix"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)


const (
	ContainerServiceSocket string = "/var/run/container_service.sock"
	DockerLauncherSocket string = "/var/run/docker_launcher.sock"
	containerPrefix string = "csaapi"
)

var path string = "http+unix://" + containerPrefix


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


func GetContainersInfo() (ContainerLists, error) {

	u := &httpunix.Transport{
		DialTimeout:           100 * time.Millisecond,
		RequestTimeout:        1 * time.Second,
		ResponseHeaderTimeout: 1 * time.Second,
	}

	u.RegisterLocation(containerPrefix, ContainerServiceSocket)

	var client = http.Client{
		Transport: u,
	}

	resp, err := client.Get(path+"/getContainersInfo");
	

	var send ContainerLists

	if err != nil {
		return send, err
	} 

	fmt.Println(resp.StatusCode)
	if resp.StatusCode == 200 {
		defer resp.Body.Close()

		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		send = ContainerLists{}
		json.Unmarshal([]byte(contents), &send)
		fmt.Println(send)

	} else {
		log.Fatal("Status : %d", resp.StatusCode)	
	}
	
	/* Stub code , it will be removed */
	send = ContainerLists{
		Cmd:           "sdfsdf",
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
	return send, nil
}
