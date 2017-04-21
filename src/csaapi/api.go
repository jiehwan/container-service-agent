package csaapi

import (
	//"bytes"
	"encoding/json"
	//"errors"
	"fmt"
	//"io/ioutil"
	"log"
	"net"
	//"net/http"
	"os"
	//	"time"
	"../types/csac"
	"../types/dockerlauncher"

	"strings"

	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	//"github.com/docker/docker/integration-cli/cli"
	//"github.com/docker/docker/api/types/container"
	//"github.com/docker/docker/integration-cli/checker"
	//icmd "github.com/docker/docker/pkg/testutil/cmd"
	//"github.com/go-check/check"
)

const (
	ContainerServiceSocket string = "/var/run/container_service.sock"
	DockerLauncherSocket   string = "/var/run/docker_launcher.sock"
)

func GetHardwareAddress() (string, error) {

	//----------------------
	// Get the local machine IP address
	// https://www.socketloop.com/tutorials/golang-how-do-I-get-the-local-ip-non-loopback-address
	//----------------------

	addrs, err := net.InterfaceAddrs()

	if err != nil {
	     fmt.Println(err)
	}

	var currentIP, currentNetworkHardwareName string

	for _, address := range addrs {

	     // check the address type and if it is not a loopback the display it
	     // = GET LOCAL IP ADDRESS
	     if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
	             if ipnet.IP.To4() != nil {
	                     fmt.Println("Current IP address : ", ipnet.IP.String())
	                     currentIP = ipnet.IP.String()
	             }
	     }
	}

	fmt.Println("------------------------------")
	fmt.Println("We want the interface name that has the current IP address")
	fmt.Println("MUST NOT be binded to 127.0.0.1 ")
	fmt.Println("------------------------------")

	// get all the system's or local machine's network interfaces

	interfaces, _ := net.Interfaces()
	for _, interf := range interfaces {

	     if addrs, err := interf.Addrs(); err == nil {
	             for index, addr := range addrs {
	                     fmt.Println("[", index, "]", interf.Name, ">", addr)

	                     // only interested in the name with current IP address
	                     if strings.Contains(addr.String(), currentIP) {
	                             fmt.Println("Use name : ", interf.Name)
	                             currentNetworkHardwareName = interf.Name
	                     }
	             }
	     }
	}

	fmt.Println("------------------------------")

	// extract the hardware information base on the interface name
	// capture above
	netInterface_, err := net.InterfaceByName(currentNetworkHardwareName)

	if err != nil {
	     fmt.Println(err)
	}

	name_ := netInterface_.Name
	macAddress_ := netInterface_.HardwareAddr

	fmt.Println("Hardware name : ", name_)
	fmt.Println("MAC address : ", macAddress_)

	// verify if the MAC address can be parsed properly
	hwAddr_, err := net.ParseMAC(macAddress_.String())

	if err != nil {
	     fmt.Println("No able to parse MAC address : ", err)
	     os.Exit(-1)
	}

	fmt.Printf("Physical hardware address : %s \n", hwAddr_.String())
	

	return macAddress_.String(), nil
}

func GetContainersInfo() (csac.ContainerLists, error) {

	var send csac.ContainerLists

	log.Printf("api.GetContainersInfo ~~~")

	// stub : "GET", "/v1/getContainersInfo"
	///////////////////////////////////////////////////////////

	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	contents, err := cli.ContainerList(context.Background(), 
									types.ContainerListOptions{All: true})
	if err != nil {
		panic(err)
	}

	// refer type def : https://github.com/docker/engine-api/blob/master/types/types.go#L148
	fmt.Printf("contents=%s\n", contents)

	
	var numOfList int = len(contents)
	log.Printf("numOfList[%d]\n", numOfList)

	send.Cmd = "GetContainersInfo"
	send.ContainerCount = numOfList

	macaddress, _ := GetHardwareAddress()

	send.DeviceID = macaddress
	log.Printf("send.DeviceID[%s]\n", send.DeviceID)

	for i := 0; i < numOfList; i++ {
		var containerValue = csac.ContainerInfo{
			ContainerName:   contents[i].Names[0],
			ImageName:       contents[i].Image,
			ContainerStatus: contents[i].State,
		}

		send.Container = append(send.Container, containerValue)
		log.Printf("[%d]-[%s]", i, send.Container)
	}

	log.Printf("[%s]", send)

	return send, nil
}

func UpdateImage(data csac.UpdateImageParams) (csac.UpdateImageReturn, error) {
	var send csac.UpdateImageReturn

	send_str, _ := json.Marshal(data)
	fmt.Println(string(send_str))

	// stub : "POST", "/v1/updateImage"
	///////////////////////////////////////////////////////////

	contents := ""

	///////////////////////////////////////////////////////////
	///////////////////////////////////////////////////////////

	object := dockerlauncher.UpdateImageReturn{}

	json.Unmarshal([]byte(contents), &object)
	log.Printf("object [%s]\n", object)

	send.Cmd = "UpdateImage"

	macaddress, _ := GetHardwareAddress()

	log.Printf("macaddress[%s]\n", macaddress)
	send.DeviceID = macaddress

	send.UpdateState = object.State.CurrentState
	log.Printf("send.UpdateState[%s]\n", send.UpdateState)

	return send, nil
}
