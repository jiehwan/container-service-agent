package main

import (
	"csaapi"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
)

const ContainerServiceSocket string = "/var/run/container_service.sock"
const DockerLauncherSocket string = "/var/run/docker_launcher.sock"

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

func getContainersInfo() (csaapi.ContainerLists, error) {
	log.Printf("getContainersInfo")

	c, err := net.Dial("unix", DockerLauncherSocket)
	if err != nil {
		log.Fatal("Dial error", err)
	}
	defer c.Close()

	var name string = "getContainersInfo"

	length := len(name)
	/*
		var size = strconv.Itoa(length)
		fmt.Printf("size: %s\n", size)
	*/
	message := make([]byte, 0, length)

	//message = append(message, size...)
	//message = append(message, " "...)

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
		//fmt.Println("buf:[%s]", buf)
		data = append(data, buf...)
	}

	//var resSize = string(data[:1])

	/*var n, err1 = strconv.Atoi(resSize)
	_ = err1

	realData := data[1 : n+1]
	fmt.Printf("%s\n", realData)*/

	//var jsonData = string(data)

	log.Printf("%s\n", data)

	//json.Unmarshal([]byte(data), &rcv)
	//fmt.Println(rcv)
	//fmt.Println(rcv.Cmd)

	// Need to parse json
	//Stub Return
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

	return send, nil
}

func GetContainersInfoHandler(writer http.ResponseWriter, request *http.Request) {
	log.Printf("Enter GetContainersInfoHandler")
	sendResponse, sendError := responseSenders(writer)
	if containersInfo, err := getContainersInfo(); err != nil {

		sendError(err)
	} else {
		//payload := make(map[string][]string)
		//payload["containersInfo"] = containersInfo

		log.Printf("%s\n", containersInfo)
		//sendResponse(payload, "", http.StatusOK)
		sendResponse("OK", "", http.StatusOK)
	}
}

func setupApi(router *mux.Router) {
	router.HandleFunc("/getContainersInfo", func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("/getContainersInfo")
		sendResponse, sendError := responseSenders(writer)
		if containersInfo, err := getContainersInfo(); err != nil {
			log.Printf("error !!!!")
			sendError(err)
		} else {
			//payload := make(map[string][]string)
			//payload["containersInfo"] = containersInfo

			log.Printf("%s\n", containersInfo)
			//sendResponse(payload, "", http.StatusOK)
			sendResponse("OK", "", http.StatusOK)
		}

		//fmt.Fprintln(writer, "OK")
		//apiv1 := router.PathPrefix("/v1").Subrouter()
		//apiv1.HandleFunc("/getContainersInfo", GetContainersInfoHandler).Methods("GET")
	})
}

func handler(writer http.ResponseWriter, request *http.Request) {
	sendResponse, sendError := responseSenders(writer)
	_ = sendResponse
	_ = sendError
	sendResponse("OK", "", http.StatusOK)
}

func main() {
	log.Printf("Container-Service Agent starting")
	listenAddress := ContainerServiceSocket
	router := mux.NewRouter()
	setupApi(router)

	/*router.HandleFunc("/", handler)

	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		t, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		fmt.Println(t)
		return nil
	})
	http.Handle("/", router)
	*/
	defer os.Remove(ContainerServiceSocket)

	if listener, err := net.Listen("unix", listenAddress); err != nil {
		log.Fatalf("Could not listen on %s: %v", listenAddress, err)
		return
	} else {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go func() {
			for sig := range c {
				// sig is a ^C, handle it
				log.Printf("Caught signal %s: shutting down.", sig)
				listener.Close()
				os.Remove(ContainerServiceSocket)
				os.Exit(0)
			}
		}()

		log.Printf("Starting HTTP server on %s\n", listenAddress)
		if err = http.Serve(listener, router); err != nil {
			log.Fatalf("Could not start HTTP server: %v", err)
		}
	}

}
