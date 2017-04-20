package csaapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tv42/httpunix"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
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

func getHardwareAddress() (string, error) {

	currentNetworkHardwareName := "eth0"
	netInterface, err := net.InterfaceByName(currentNetworkHardwareName)

	if err != nil {
		fmt.Println(err)
	}

	name := netInterface.Name
	macAddress := netInterface.HardwareAddr

	log.Printf("Hardware name : %s\n", string(name))
	log.Printf("MAC address : %s\n", string(macAddress))

	// verify if the MAC address can be parsed properly
	hwAddr, err := net.ParseMAC(macAddress.String())

	if err != nil {
		log.Printf("No able to parse MAC address : %s\n", err)
		os.Exit(-1)
	}

	log.Printf("Physical hardware address : %s \n", hwAddr.String())

	return hwAddr.String(), nil
}

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
	log.Printf("csaapi : %d", resp.StatusCode)
	var send csac.ContainerLists

	if err != nil {
		log.Printf("err [%s]", err)
		return send, err
	}

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
		log.Printf("numOfList[%d]\n", numOfList)

		send.Cmd = "GetContainersInfo"
		send.ContainerCount = numOfList

		macaddress, err := getHardwareAddress()

		log.Printf("macaddress[%s]\n", macaddress)
		send.DeviceID = macaddress

		for i := 0; i < numOfList; i++ {
			var containerValue = csac.ContainerInfo{
				ContainerName:   lists.Containers[i].ContainerName,
				ImageName:       lists.Containers[i].ImageName,
				ContainerStatus: lists.Containers[i].ContainerStatus,
			}

			send.Container = append(send.Container, containerValue)
			log.Printf("[%d]-[%s]", i, send.Container)
		}
		log.Printf("[%s]", send)

	} else {
		log.Printf("Status : %d", resp.StatusCode)

		return send, errors.New("Not able to use api")
	}

	return send, nil
}

func UpdateImage(csac.UpdateImageParams) (csac.UpdateImageReturn, error) {

	var send csac.UpdateImageReturn

	send.Cmd = "UpdateImage"
	macaddress, err := getHardwareAddress()
	if err != nil {
		return send, err
	}

	log.Printf("macaddress[%s]\n", macaddress)
	send.DeviceID = macaddress

	send.DeviceID = "ARTIK710-1"
	send.UpdateState = "Started"

	return send, nil
}
