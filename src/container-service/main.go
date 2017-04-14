package main

import (
	"fmt"
	"log"
	
	"golang.org/x/net/websocket"
	"os"
	"io"
	"strings"
	"net"
	"net/http"
	"net/url"
	"net/http/httputil"

	"encoding/json"	
)

var wss_server_url = "ws://10.113.62.204:4000"
var wss_server_origin = "ws://10.113.62.204:4000"

type Command struct {
	Cmd string `json:"cmd"`
}


func main() {

	ws, err := ProxyDial(wss_server_url, "tcp", wss_server_origin)

	if err != nil {
		log.Fatal("ProxyDial : ", err)
	}

	defer ws.Close()

	/* connect test2 : message driven
	*/
	messages := make(chan string)
	go wsReceive(ws, messages)

	name, _ := os.Hostname()
    err = wsReqeustConnection(ws, name)

    for{
    	msg := <-messages

    	rcv := Command{}
		json.Unmarshal([]byte(msg), &rcv)
	    fmt.Println(rcv.Cmd)

	    switch (rcv.Cmd) {
    	case "connected" :
    		log.Printf("connected succefully~~")

		case "getContainerLists" :
    		log.Printf("command <getContainerLists>")
    		wsSendContainerLists(ws)

	    default :
	    	log.Printf("add command of {%s}", rcv.Cmd)
	    }


    }
}

func wsReceive(ws *websocket.Conn, chan_msg chan string) (err error) {
	var read_buf string

	for {
		err = websocket.Message.Receive(ws, &read_buf)
		if (err != nil) {
			log.Fatal(err)
		}
		log.Printf("received: %s", read_buf)
		chan_msg <- read_buf
	}
	return err
}

type ContainerInfo struct {
    ContainerID string `json:"container_id"`
	ContainerStatus string `json:"container_status"`
}
type ContainerLists struct {
	Cmd string `json:"cmd"`
	ContainerCount int `json:"container_count"`
	Container []ContainerInfo `json:"container"`
}

func wsSendContainerLists(ws *websocket.Conn) (err error) {

	//First.. OK
	send := ContainerLists{
		Cmd : "getContainerLists",
		ContainerCount : 2,
		Container :[]ContainerInfo{
			{ 
				ContainerID : "1111",
				ContainerStatus : "running",
			},
			{
				ContainerID : "2222",
				ContainerStatus : "exited",
			},
		},
	}

	websocket.JSON.Send(ws, send)

	return nil
}


func wsTest1(ws *websocket.Conn) (err error){
	name, _ := os.Hostname()
    err = wsReqeustConnection(ws, name)

    // receive connection token
    Token, err := wsReceiveConnection(ws)
	log.Printf("recv.Token = '%s'", Token)

	return err
}


type ConnectReq struct {
	Cmd string `json:"cmd"`
	Name string `json:"name"`
}

func wsReqeustConnection(ws *websocket.Conn, name string) (err error) {
	send := ConnectReq{}
    send.Cmd = "request"
    send.Name = name

	websocket.JSON.Send(ws, send)

	return nil
}


type ConnectedResp struct {
	Cmd string `json:"cmd"`
	Token string `json:"token"`
	Clinetnum int `json:"clientnum"`
}

func wsReceiveConnection(ws *websocket.Conn) (Token string, err error) {
	recv := ConnectedResp{}

	err = websocket.JSON.Receive(ws, &recv)
	if(err != nil) {
		log.Fatal(err)
	}

	return recv.Token, err
}


func ProxyDial(url_, protocol, origin string) (ws *websocket.Conn, err error) {

	log.Printf("http_proxy {%s}\n", os.Getenv("HTTP_PROXY"))

	// comment out in case of testing without proxy
	if strings.Contains(url_, "10.113.") {
		return websocket.Dial(url_, protocol, origin)
	}

	if os.Getenv("HTTP_PROXY") == "" {
		return websocket.Dial(url_, protocol, origin)
	}

	purl, err := url.Parse(os.Getenv("HTTP_PROXY"))
	if err != nil {
		log.Fatal("Parse : ", err)
		return nil, err
	}

	log.Printf("====================================")
	log.Printf("    websocket.NewConfig")
	log.Printf("====================================")
	config, err := websocket.NewConfig(url_, origin)
	if err != nil {
		log.Fatal("NewConfig : ", err)
		return nil, err
	}

	if protocol != "" {
		config.Protocol = []string{protocol}
	}

	log.Printf("====================================")
	log.Printf("    HttpConnect")
	log.Printf("====================================")
	client, err := HttpConnect(purl.Host, url_)
	if err != nil {
		log.Fatal("HttpConnect : ", err)
		return nil, err
	}

	log.Printf("====================================")
	log.Printf("    websocket.NewClient")
	log.Printf("====================================")
	return websocket.NewClient(config, client)
}


func HttpConnect(proxy, url_ string) (io.ReadWriteCloser, error) {
	log.Printf("proxy =", proxy)
	proxy_tcp_conn, err := net.Dial("tcp", proxy)
	if err != nil {
		return nil, err
	}
	log.Printf("proxy_tcp_conn =", proxy_tcp_conn)
	log.Printf("url_ =", url_)

	turl, err := url.Parse(url_)
	if err != nil {
		log.Fatal("Parse : ", err)
		return nil, err
	}
	
	log.Printf("proxy turl.Host =", string(turl.Host))


	req := http.Request{
		Method: "CONNECT",
		URL:    &url.URL{},
		Host:   turl.Host,
	}

	/*
	// origin
	req := http.Request{
		Method: "CONNECT",
		URL:    &url.URL{},
		Host:   turl.Host,
	}
	*/

	proxy_http_conn := httputil.NewProxyClientConn(proxy_tcp_conn, nil)
	//cc := http.NewClientConn(proxy_tcp_conn, nil)

	log.Printf("proxy_http_conn =", proxy_http_conn)	

	resp, err := proxy_http_conn.Do(&req)
	if err != nil && err != httputil.ErrPersistEOF {
		log.Fatal("ErrPersistEOF : ", err)
		return nil, err
	}
	log.Printf("proxy_http_conn<resp> =", (resp))

	rwc, _ := proxy_http_conn.Hijack()

	return rwc, nil
	
}


// return Handler (A Handler reponds to an HTTP request)
func websocketProxy(target string) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                d, err := net.Dial("tcp", target)
                if err != nil {
                        http.Error(w, "Error contacting backend server.", 500)
                        log.Printf("Error dialing websocket backend %s: %v", target, err)
                        return
                }
                hj, ok := w.(http.Hijacker)
                if !ok {
                        http.Error(w, "Not a hijacker?", 500)
                        return
                }
                nc, _, err := hj.Hijack()
                if err != nil {
                        log.Printf("Hijack error: %v", err)
                        return
                }
                defer nc.Close()
                defer d.Close()

                err = r.Write(d)
                if err != nil {
                        log.Printf("Error copying request to target: %v", err)
                        return
                }

                errc := make(chan error, 2)
                cp := func(dst io.Writer, src io.Reader) {
                        _, err := io.Copy(dst, src)
                        errc <- err
                }
                go cp(d, nc)
                go cp(nc, d)
                <-errc
        })
    }

func json_marshal() {
	// convert from struct to string
    send := ConnectedResp{}
    send.Cmd = "request"
    send.Token = "1234"
    send.Clinetnum = 88

    send_str, _ := json.Marshal(send)
    fmt.Println(string(send_str))
}

func json_unmarshal() {
	// convert from string to struct
	rcv_str := `{"cmd": "connected" 
			, "token": "test-token"
			, "clinetnum": 3}`
	rcv := ConnectedResp{}
	json.Unmarshal([]byte(rcv_str), &rcv)
    fmt.Println(rcv)
    fmt.Println(rcv.Cmd)
}

