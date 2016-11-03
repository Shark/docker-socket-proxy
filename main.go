package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func handleErr(writer http.ResponseWriter, status int, err error) {
	msg := fmt.Sprintf("%v", err)
	log.Print(msg)

	msgJson := struct {
		Message string `json:"message"`
	}{msg}
	bytes, err := json.Marshal(msgJson)
	if err != nil {
		log.Printf("Error encoding json: %v", err)
		return
	}

	writer.WriteHeader(status)
	writer.Write(bytes)
}

type handler struct {
	outsock string
}

func (h *handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		handleErr(writer, 403, fmt.Errorf("Unauthorized request to %s %s", request.Method, request.URL.Path))
		return
	}

	conn, err := net.DialUnix("unix", nil, &net.UnixAddr{h.outsock, "unix"})
	if err != nil {
		handleErr(writer, 500, err)
		return
	}
	defer conn.Close()
	err = request.Write(conn)
	if err != nil {
		handleErr(writer, 500, err)
		return
	}
	resp, err := http.ReadResponse(bufio.NewReader(conn), request)
	if err != nil {
		handleErr(writer, 500, err)
		return
	}
	defer resp.Body.Close()
	for k, v := range resp.Header {
		writer.Header()[k] = v
	}
	writer.WriteHeader(resp.StatusCode)

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')

		if err == io.EOF {
			return
		} else if err != nil {
			log.Fatalf("Error reading body: %v", err)
			return
		}

		writer.Write(line)
		if flusher, ok := writer.(http.Flusher); ok {
			flusher.Flush()
		}
	}
}

type config struct {
	insock  string
	outsock string
}

func parseConfig() config {
	insock := flag.String("in", "docker-socket-proxy.sock", "Incoming socket")
	outsock := flag.String("out", "/var/run/docker.sock", "Outgoing socket (i.e. Docker socket)")
	flag.Parse()

	return config{*insock, *outsock}
}

func main() {
	config := parseConfig()

	sock, err := net.Listen("unix", config.insock)
	if err != nil {
		log.Fatalf("Can not listen: %v", err)
		return
	}

	myHandler := &handler{config.outsock}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func(c chan os.Signal) {
		sig := <-sigChan
		log.Printf("Got signal %v: exiting", sig)
		sock.Close()
		os.Exit(0)
	}(sigChan)

	err = http.Serve(sock, myHandler)
	if err != nil {
		log.Fatalf("Failed to serve http: %v", err)
		return
	}
}
