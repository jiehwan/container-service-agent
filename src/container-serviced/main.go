package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"csaapi"
	"fmt"
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

func getContainersInfo() ([]byte, error) {
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
	message := make([]byte, 0, length)
	message = append(message, name...)

	log.Printf("before send: %s\n", string(message))

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
		buf := make([]byte, 5465)
		nr, err := c.Read(buf)
		if err != nil {
			break
		}

		buf = buf[:nr]
		data = append(data, buf...)
	}

	log.Printf("%s\n", data)
	// Need to parse json
	//Stub Return
	/*send = csaapi.ContainerLists{
		Cmd:            "getContainerLists",
		ContainerCount: 2,
		Container: []csaapi.ContainerInfo{
			{
				ContainerID:     "api-1111",
				ContainerStatus: "running",
			},
			{
				ContainerID:     "api-2222",
				ContainerStatus: "exited",
			},
		},
	}*/

	
	send := csaapi.ContainerLists{
		Cmd:            "getContainerLists",
		ContainerCount: 2,
		Container: []csaapi.ContainerInfo{
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

	send_str, _ = json.Marshal(send)
	fmt.Println(string(send_str))


	return send_str, nil
}

func GetContainersInfoHandler(writer http.ResponseWriter, request *http.Request) {
	log.Printf("Enter GetContainersInfoHandler")

	sendResponse, sendError := responseSenders(writer)
	if containersInfo, err := getContainersInfo(); err != nil {
		sendError(err)
	} else {
		log.Printf("%s\n", containersInfo)
		sendResponse("OK", "", http.StatusOK)
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
