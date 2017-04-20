package main

import (
	"bytes"
	"csaapi"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"types/csac"
	"types/dockerlauncher"
)

const (
	MaxCommandLength int = 10
)

// APIResponse The api response sent from go supervisor
type APIResponse struct {
	Data  interface{}
	Error string
}

func jsonResponse(writer http.ResponseWriter, response interface{}, status int) {
	jsonBody, err := json.Marshal(response)
	if err != nil {
		log.Printf("Could not marshal JSON for %+v\n", response)
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	writer.Write(jsonBody)
}

func responseSenders(writer http.ResponseWriter) (sendResponse func(interface{}, string, int), sendError func(error)) {
	sendResponse = func(data interface{}, errorMsg string, statusCode int) {
		jsonResponse(writer, APIResponse{data, errorMsg}, statusCode)
	}
	sendError = func(err error) {
		sendResponse("Error", err.Error(), http.StatusInternalServerError)
	}
	return
}

func getDockerLauncherInfo_Stub() dockerlauncher.GetContainersInfoReturn {
	send := dockerlauncher.GetContainersInfoReturn{
		Containers: []dockerlauncher.Container{
			{
				ContainerName:   "aaaa",
				ImageName:       "tizen1",
				ContainerStatus: "created",
			},
			{
				ContainerName:   "bbbb",
				ImageName:       "tizen2",
				ContainerStatus: "exited",
			},
		},
	}

	return send
}

func updateImage_Stub() dockerlauncher.UpdateImageReturn {
	send := dockerlauncher.UpdateImageReturn{
		State: dockerlauncher.DeviceState{
			CurrentState: "Updating",
		},
	}

	return send
}

func getContainersInfo() ([]byte, error) {
	log.Printf("getContainersInfo")

	/*stub := getDockerLauncherInfo_Stub()
	var send_stub []byte

	send_stub, _ = json.Marshal(stub)
	log.Printf(string(send_stub))

	return send_stub, nil
	*/
	var send_str []byte
	c, err := net.Dial("unix", csaapi.DockerLauncherSocket)
	if err != nil {
		log.Fatal("Dial error", err)
		return send_str, nil
	}

	defer c.Close()

	send := dockerlauncher.Cmd{}
	send.Command = "GetContainersInfo"

	send_str, _ = json.Marshal(send)
	log.Printf(string(send_str))

	length := len(send_str)

	message := make([]byte, 0, length)
	message = append(message, send_str...)

	_, err = c.Write([]byte(message))
	if err != nil {
		log.Printf("error: %v\n", err)
	}

	log.Printf("sent: %s\n", message)
	err = c.(*net.UnixConn).CloseWrite()
	if err != nil {
		log.Printf("error: %v\n", err)

	}

	data := make([]byte, 0)
	for {
		dataBuf := make([]byte, 1024)
		nr, err := c.Read(dataBuf)
		if err != nil {
			break
		}

		log.Printf("nr size [%d]", nr)
		if nr == 0 {
			break
		}

		dataBuf = dataBuf[:nr]
		data = append(data, dataBuf...)
	}
	log.Printf("receive data[%s]\n", string(data))
	//delete null character
	withoutNull := bytes.Trim(data, "\x00")

	rcv := dockerlauncher.Cmd{}
	err = json.Unmarshal([]byte(withoutNull), &rcv)
	log.Printf("rcv.Command = %s", rcv.Command)

	if rcv.Command == "GetContainersInfo" {
		log.Printf("Success\n")
		return withoutNull, nil
	} else {
		log.Printf("error commnad[%s]\n", err)
	}

	log.Printf("end\n")
	return send_str, nil
}

func updateImageRequest(ImageName, ContainerName string) ([]byte, error) {
	log.Printf("updateImageRequest")

	/*stub := updateImage_Stub()
	var send_stub []byte

	send_stub, _ = json.Marshal(stub)
	log.Printf(string(send_stub))

	return send_stub, nil
	*/
	var send_str []byte
	c, err := net.Dial("unix", csaapi.DockerLauncherSocket)
	if err != nil {
		log.Fatal("Dial error", err)
		return send_str, nil
	}

	defer c.Close()

	send := dockerlauncher.UpdateImageParameters{}
	send.Command = "UpdateImage"

	send.Param = dockerlauncher.UpdateParam{
		ContainerName: ContainerName,
		ImageName:     ImageName,
	}

	send_str, _ = json.Marshal(send)
	log.Printf(string(send_str))

	length := len(send_str)

	message := make([]byte, 0, length)
	message = append(message, send_str...)

	_, err = c.Write([]byte(message))
	if err != nil {
		log.Printf("error: %v\n", err)
	}

	log.Printf("sent: %s\n", message)
	err = c.(*net.UnixConn).CloseWrite()
	if err != nil {
		log.Printf("error: %v\n", err)

	}

	data := make([]byte, 0)
	for {
		dataBuf := make([]byte, 1024)
		nr, err := c.Read(dataBuf)
		if err != nil {
			break
		}

		log.Printf("nr size [%d]", nr)
		if nr == 0 {
			break
		}

		dataBuf = dataBuf[:nr]
		data = append(data, dataBuf...)
	}
	log.Printf("receive data[%s]\n", string(data))
	//delete null character
	withoutNull := bytes.Trim(data, "\x00")

	rcv := dockerlauncher.Cmd{}
	err = json.Unmarshal([]byte(withoutNull), &rcv)
	log.Printf("rcv.Command = %s", rcv.Command)

	if rcv.Command == "UpdateImage" {
		log.Printf("Success\n")
		return withoutNull, nil
	} else {
		log.Printf("error commnad[%s]\n", err)
	}

	log.Printf("end\n")
	return send_str, nil
}

func parseUpdateImageParam(request *http.Request) (ImageName, ContainerName string, err error) {

	var body csac.UpdateImageParams

	decoder := json.NewDecoder(request.Body)
	decoder.Decode(&body)

	log.Printf("body.ImageName = %s", body.ImageName)
	log.Printf("body.ContainerName = %s", body.ContainerName)

	ImageName = body.ImageName
	ContainerName = body.ContainerName

	return ImageName, ContainerName, err
}

func GetContainersInfoHandler(writer http.ResponseWriter, request *http.Request) {
	log.Printf("Enter GetContainersInfoHandler")

	if containersInfo, err := getContainersInfo(); err != nil {
		log.Printf("Error GetContainersInfoHandler[%s]", err)
		writer.WriteHeader(http.StatusInternalServerError)
	} else {
		log.Printf("Success GetContainersInfoHandler")
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		writer.Write(containersInfo)
	}
}

func UpdateImageHandler(writer http.ResponseWriter, request *http.Request) {
	log.Printf("Enter UpdateImageHandler")

	imageName, containerName, err := parseUpdateImageParam(request)
	if err != nil {
		log.Printf("Error here [%s]", err)
	}
	if updateImageState, err := updateImageRequest(imageName, containerName); err != nil {
		log.Printf("Error UpdateImageHandler[%s]", err)
		writer.WriteHeader(http.StatusInternalServerError)
	} else {
		log.Printf("Success UpdateImageHandler")
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		writer.Write(updateImageState)
	}

}

func htmlHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(writer, "OK")
}

func setupApi(r *mux.Router) {
	r.PathPrefix("/").
		Path("/ping").
		Methods("GET").
		HandlerFunc(htmlHandler)

	s := r.PathPrefix("/v1").Subrouter()
	s.HandleFunc("/getContainersInfo", GetContainersInfoHandler).Methods("GET")
	s.HandleFunc("/updateImage", UpdateImageHandler).Methods("POST")
}

func main() {
	log.Printf("Container-Service Agent starting")
	listenAddress := csaapi.ContainerServiceSocket
	router := mux.NewRouter()
	setupApi(router)

	if listener, err := net.Listen("unix", listenAddress); err != nil {
		log.Fatalf("Could not listen on %s: %v", listenAddress, err)
		return
	} else {

		defer listener.Close()

		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGKILL)
		go func(c chan os.Signal) {
			sig := <-c
			log.Printf("Caught signal %s: shutting down.", sig)
			listener.Close()
			os.Exit(0)
		}(sigc)

		log.Printf("Starting HTTP server on %s\n", listenAddress)
		if err = http.Serve(listener, router); err != nil {
			log.Fatalf("Could not start HTTP server: %v", err)
		}
	}

}
