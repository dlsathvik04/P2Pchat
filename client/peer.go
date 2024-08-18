package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"
)

type Peer struct {
	StunUrl          string
	Port             int
	Username         string
	password         string
	handleConnection func(net.Conn)
}

func (p *Peer) Dial(username string) (net.Conn, error) {
	stunUrl := p.StunUrl
	dialer := &net.Dialer{
		Timeout: 3 * time.Second,
	}
	transport := &http.Transport{
		DisableKeepAlives: true,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			conn, err := dialer.DialContext(ctx, network, addr)
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
	}

	httpClient := &http.Client{
		Transport: transport,
	}

	reqBody, _ := json.Marshal(map[string]string{
		"username": p.Username,
		"password": p.password,
	})
	authReq, err := http.NewRequest("POST", stunUrl+"/login", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}
	authResp, err := httpClient.Do(authReq)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil, err
	}
	var response struct{ Token string }
	json.NewDecoder(authResp.Body).Decode(&response)
	fmt.Println(response)
	authResp.Body.Close()

	queryReq, err := http.NewRequest("GET", stunUrl+"/getip?username=dlsathvik04", nil)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil, err
	}
	queryReq.Header.Add("authtoken", response.Token)
	queryResp, err := httpClient.Do(queryReq)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil, err
	}
	var queryResBody struct {
		Username string
		Addr     string
	}
	json.NewDecoder(queryResp.Body).Decode(&queryResBody)
	fmt.Println(queryResBody)
	queryResp.Body.Close()
	//TODO - implement tcp dialing

	conn, err := net.Dial("tcp", queryResBody.Addr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (p *Peer) StartAccepting() {
	stunUrl := p.StunUrl
	var localPort int
	dialer := &net.Dialer{
		Timeout: 3 * time.Second,
		LocalAddr: &net.TCPAddr{
			Port: p.Port,
		},
	}
	transport := &http.Transport{
		DisableKeepAlives: true,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			conn, err := dialer.DialContext(ctx, network, addr)
			if err != nil {
				return nil, err
			}
			localAddr := conn.LocalAddr().(*net.TCPAddr)
			localPort = localAddr.Port

			return conn, nil
		},
	}

	httpClient := &http.Client{
		Transport: transport,
	}

	reqBody, _ := json.Marshal(map[string]string{
		"username": p.Username,
		"password": p.password,
	})
	authReq, err := http.NewRequest("POST", stunUrl+"/login", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	authResp, err := httpClient.Do(authReq)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	var response struct{ Token string }
	json.NewDecoder(authResp.Body).Decode(&response)
	fmt.Println(response)
	authResp.Body.Close()

	stunReq, err := http.NewRequest("GET", stunUrl+"/stun?register=true", nil)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	stunReq.Header.Add("authtoken", response.Token)
	stunResp, err := httpClient.Do(stunReq)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	var stunresbody struct {
		Username string
		Addr     string
	}
	json.NewDecoder(stunResp.Body).Decode(&stunresbody)
	fmt.Println(stunresbody)
	stunResp.Body.Close()
	fmt.Println("Local port used:", localPort)
	fmt.Println("Response status code:", authResp.StatusCode)

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(localPort))
	if err != nil {
		fmt.Println("Error starting TCP listener:", err)
		return
	}
	defer listener.Close()

	fmt.Println("TCP listener started on port", localPort)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		fmt.Println("Connection accepted from", conn.RemoteAddr())
		go p.handleConnection(conn)
	}
}
