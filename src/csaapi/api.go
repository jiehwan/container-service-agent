package csaapi

import (
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/tv42/httpunix"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	ContainerServiceSocket string = "/var/run/container_service.sock"
	DockerLauncherSocket   string = "/var/run/docker_launcher.sock"
	containerPrefix        string = "csaapi"
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

	resp, err := client.Get(path + "/v1/getContainersInfo")

	var send ContainerLists

	if err != nil {
		return send, err
	}

	log.Printf("csaapi : %d", resp.StatusCode)

	if resp.StatusCode == 200 {
		defer resp.Body.Close()

		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		lists := make([]types.Container, 0)

		json.Unmarshal([]byte(contents), &lists)
		log.Printf("List [%s]", lists)

		send = ContainerLists{
			Cmd:            lists[0].Command,
			ContainerCount: 1,
			Container: []ContainerInfo{
				{
					ContainerID:     lists[0].ID,
					ContainerStatus: lists[0].Status,
				},
			},
		}

	} else {
		log.Printf("Status : %d", resp.StatusCode)
	}

	return send, nil
}
