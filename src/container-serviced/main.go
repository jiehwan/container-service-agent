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
	"strconv"
	"syscall"
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
		Command: "GetContainersInfo",

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

func getContainersInfo() ([]byte, error) {
	log.Printf("getContainersInfo")

	/*
		stub := getDockerLauncherInfo_Stub()
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

func getContainersInfo2() ([]byte, error) {
	log.Printf("getContainersInfo")

	var send_str []byte
	c, err := net.Dial("unix", csaapi.DockerLauncherSocket)
	if err != nil {
		log.Fatal("Dial error", err)
		return send_str, nil
	}

	defer c.Close()

	var name string = "getContainersInfo"
	length := len(name)

	command_size := strconv.Itoa(length)
	command_size_len := len(command_size)
	blank := " "
	blank_len := len(" ")

	message := make([]byte, 0, length+command_size_len+blank_len)
	message = append(message, command_size...)
	message = append(message, blank...)
	message = append(message, name...)

	_, err = c.Write([]byte(message))
	if err != nil {
		log.Printf("error: %v\n", err)
	}

	log.Printf("sent: %s\n", message)
	err = c.(*net.UnixConn).CloseWrite()
	if err != nil {
		log.Printf("error: %v\n", err)

	}

	// Wating message to find size
	// buf is size
	buf := make([]byte, MaxCommandLength)
	for {
		nr, _ := c.Read(buf)
		if nr != 0 {
			break
		}
	}

	var count int
	for i, v := range buf {
		if v == ' ' {
			count = i
			log.Printf("Position of blank[%d]\n", i)
			break
		}
	}

	sizeArray := make([]byte, count)
	sizeArray = buf[0:count]

	log.Printf("sizearray [%s]\n", string(sizeArray))
	num, _ := strconv.Atoi(string(sizeArray))

	fmt.Println("Total JSON size is ", num)

	// wating real message
	data := make([]byte, num)
	data = append(data, buf[count+1:MaxCommandLength]...)
	log.Printf("before reading [%s]\n", string(data))

	var checkReceiveSize int = MaxCommandLength - count - 1

	for {
		dataBuf := make([]byte, num-checkReceiveSize)
		nr, err := c.Read(dataBuf)
		if err != nil {
			break
		}

		fmt.Printf("receive data[%s]\n", string(dataBuf))

		dataBuf = dataBuf[:nr]
		checkReceiveSize += nr
		log.Printf("CheckReceiveSize [%d]\n", checkReceiveSize)
		data = append(data, dataBuf...)

		if checkReceiveSize >= num {
			break
		}
	}
	//delete null character
	withoutNull := bytes.Trim(data, "\x00")

	return withoutNull, nil
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
