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
		var numOfList int = len(lists)
		log.Printf("numOfList[%d]", numOfList)

		send.Cmd = "/bin/bash"
		send.ContainerCount = numOfList

		for i := 0; i < numOfList; i++ {
			var containerValue = ContainerInfo{
				ContainerID:     lists[i].ID,
				ContainerStatus: lists[i].Status,
			}

			send.Container = append(send.Container, containerValue)
			log.Printf("[%d]-[%s]", i, send.Container)
		}

	} else {
		log.Printf("Status : %d", resp.StatusCode)
	}

	log.Printf("[%s]", send)

	return send, nil
}
