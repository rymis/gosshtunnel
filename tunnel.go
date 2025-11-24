package main

import (
	"crypto/sha3"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

type SshClient struct {
	conn *ssh.Client
	pkBlueprint string
}

func NewSshTunnel(user, host, key, blueprint string) (*SshClient, error) {
	// Try to read and parse the key:
	keyData, err := os.ReadFile(key)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(keyData)
	if err != nil {
		return nil, err
	}

	// Add default port if not specified:
	if strings.Index(host, ":") < 0 {
		host = host + ":22"
	}

	// Let's support only one algorithm taken from the key:
	alg := signer.PublicKey().Type()

	// We need to create the result here to use it's methods:
	res := &SshClient{}
	res.pkBlueprint = blueprint

	// Prepare configuration:
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback:   res.checkHostBlueprint,
		HostKeyAlgorithms: []string{alg},
	}

	// Connect to the server:
	conn, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		return nil, err
	}

	res.conn = conn

	return res, nil
}

// Check host public key blueprint:
func (sc *SshClient) checkHostBlueprint(hostname string, remote net.Addr, key ssh.PublicKey) error {
	data := key.Marshal()
	hash := sha3.Sum256(data)
	shash := base64.StdEncoding.EncodeToString(hash[:])[:8]

	if sc.pkBlueprint != "" {
		if shash != sc.pkBlueprint {
			log.Printf("ERROR: Invalid key blueprint: %s", shash)
			return fmt.Errorf("ERROR: Invalid key blueprint: %s", shash)
		}
	} else {
		sc.pkBlueprint = shash
	}

	return nil
}

// Close the connection
func (sc *SshClient) Close() error {
	return sc.conn.Close()
}

// Wait for connection to close and return the reason
func (sc *SshClient) Wait() error {
	return sc.conn.Wait()
}

// Get current public key blueprint
func (sc *SshClient) KeyBlueprint() string {
	return sc.pkBlueprint
}

// Redirect one addr:port to specific address
func (sc *SshClient) Redirect(listen, connect string) error {
	l, err := sc.conn.Listen("tcp", listen)
	if err != nil {
		return err
	}

	go func() {
		defer l.Close()

		for {
			srv, err := l.Accept()
			if err != nil {
				log.Printf("ERROR: accept from %s failed: %s", listen, err)
				break
			}

			cli, err := net.Dial("tcp", connect)
			if err != nil {
				log.Printf("ERROR: connect to %s failed: %s", connect, err)
				srv.Close()
				continue
			}

			go proxyFrom(cli, srv)
			go proxyFrom(srv, cli)
		}
	}()

	return nil
}

func proxyFrom(src, dst net.Conn) {
	var buf [8192]byte
	defer src.Close()
	defer dst.Close()

	for {
		n, err := src.Read(buf[:])
		if err != nil {
			log.Printf("Read failed with error %s", err)
			return
		}

		n, err = dst.Write(buf[:n])
		if err != nil {
			log.Printf("Write failed with error %s", err)
			return
		}
	}
}
