package csaapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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
)

const defaultTimeout = 30 * time.Second

type CSAClient struct {
	Path       string
	HTTPClient *http.Client
}

func newHTTPClient(path string, timeout time.Duration) *http.Client {
	httpTransport := &http.Transport{}

	socketPath := path
	unixDial := func(proto, addr string) (net.Conn, error) {
		return net.DialTimeout("unix", socketPath, timeout)
	}
	httpTransport.Dial = unixDial

	return &http.Client{Transport: httpTransport}
}

func NewCSAClient() (*CSAClient, error) {

	httpClient := newHTTPClient(ContainerServiceSocket, time.Duration(defaultTimeout))
	return &CSAClient{ContainerServiceSocket, httpClient}, nil
}

func (client *CSAClient) doRequest(method string, path string, body string) ([]byte, error) {
	log.Printf("doRequest Method[%s] path[%s]", method, path)

	var resp *http.Response
	var err error

	switch method {
	case "GET":
		resp, err = client.HTTPClient.Get("http://unix" + path)
	case "POST":
		reqBody := bytes.NewBufferString(body)
		log.Printf("reqBody : [%s]\n", reqBody)
		resp, err = client.HTTPClient.Post("http://unix"+path, "text/plain", reqBody)
	default:
		return nil, errors.New("Invaild Method")
	}

	if resp.StatusCode == 200 {
		defer resp.Body.Close()
		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		return contents, err
	} else {
		log.Printf("Error  : [%d]\n", resp.StatusCode)
		return nil, errors.New(string(resp.StatusCode))
	}

	return nil, err
}

func GetHardwareAddress() (string, error) {

	currentNetworkHardwareName := "eth0"
	netInterface, err := net.InterfaceByName(currentNetworkHardwareName)

	if err != nil {
		fmt.Println(err)
	}

	name := netInterface.Name
	macAddress := netInterface.HardwareAddr

	log.Printf("Hardware name : %s\n", string(name))

	hwAddr, err := net.ParseMAC(macAddress.String())

	if err != nil {
		log.Printf("No able to parse MAC address : %s\n", err)
		os.Exit(-1)
	}

	log.Printf("Physical hardware address : %s \n", hwAddr.String())

	return hwAddr.String(), nil
}

func (client *CSAClient) GetContainersInfo() (csac.ContainerLists, error) {

	var send csac.ContainerLists

	contents, err := client.doRequest("GET", "/v1/getContainersInfo", "")

	if err != nil {
		log.Printf("error [%s]", err)
		return send, err
	}

	lists := dockerlauncher.GetContainersInfoReturn{}

	json.Unmarshal([]byte(contents), &lists)
	var numOfList int = len(lists.Containers)
	log.Printf("numOfList[%d]\n", numOfList)

	send.Cmd = "GetContainersInfo"
	send.ContainerCount = numOfList

	macaddress, err := GetHardwareAddress()

	send.DeviceID = macaddress
	log.Printf("send.DeviceID[%s]\n", send.DeviceID)

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

	return send, nil
}

func (client *CSAClient) UpdateImage(data csac.UpdateImageParams) (csac.UpdateImageReturn, error) {
	var send csac.UpdateImageReturn

	send_str, _ := json.Marshal(data)
	fmt.Println(string(send_str))

	contents, err := client.doRequest("POST", "/v1/updateImage", string(send_str))

	if err != nil {
		log.Printf("error [%s]", err)
		return send, err
	}

	object := dockerlauncher.UpdateImageReturn{}

	json.Unmarshal([]byte(contents), &object)
	log.Printf("object [%s]\n", object)

	send.Cmd = "UpdateImage"

	macaddress, err := GetHardwareAddress()

	log.Printf("macaddress[%s]\n", macaddress)
	send.DeviceID = macaddress

	send.UpdateState = object.State.CurrentState
	log.Printf("send.UpdateState[%s]\n", send.UpdateState)

	return send, nil
}
