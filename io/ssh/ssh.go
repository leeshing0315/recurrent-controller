package ssh

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"net/http"
)

type SshSession struct {
	ipaddr string
	user string
	password string
}

// https://pkg.go.dev/golang.org/x/crypto@v0.0.0-20201221181555-eec23a3978ad/ssh#section-documentation
func NewConnection(ipaddr string, user string, password string)  {
	var hostKey ssh.PublicKey
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.FixedHostKey(hostKey),
	}
	// Dial your ssh server.
	conn, err := ssh.Dial("tcp", "localhost:22", config)
	if err != nil {
		log.Fatal("unable to connect: ", err)
	}
	defer conn.Close()

	// Request the remote side to open port 8080 on all interfaces.
	l, err := conn.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatal("unable to register tcp forward: ", err)
	}
	defer l.Close()

	// Serve HTTP with your SSH server acting as a reverse proxy.
	http.Serve(l, http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(resp, "Hello world!\n")
	}))
}