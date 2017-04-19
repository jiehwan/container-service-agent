package csaapi

import (
	"encoding/json"
	"github.com/tv42/httpunix"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"types/csac"
	"types/dockerlauncher"
)

const (
	ContainerServiceSocket string = "/var/run/container_service.sock"
	DockerLauncherSocket   string = "/var/run/docker_launcher.sock"
	containerPrefix        string = "csaapi"
)

var path string = "http+unix://" + containerPrefix

func GetContainersInfo() (csac.ContainerLists, error) {

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

	var send csac.ContainerLists

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

		lists := dockerlauncher.GetContainersInfoReturn{}

		json.Unmarshal([]byte(contents), &lists)
		log.Printf("List [%s]", lists)
		var numOfList int = len(lists.Containers)
		log.Printf("numOfList[%d]", numOfList)

		send.Cmd = "GetContainersInfo"
		send.ContainerCount = numOfList
		send.DeviceID = "docker for Tizen"
		for i := 0; i < numOfList; i++ {
			var containerValue = csac.ContainerInfo{
				ContainerName:   lists.Containers[i].ContainerID,
				ImageName:       lists.Containers[i].ImageName,
				ContainerStatus: lists.Containers[i].ContainerStatus,
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

func UpdateImage(csac.UpdateImageParams) (csac.UpdateImageReturn, error) {

	var send csac.UpdateImageReturn

	send.Cmd = "UpdateImage"
	send.DeviceID = "docker for Tizen"
	send.UpdateState = "Started"

	return send, nil
}
