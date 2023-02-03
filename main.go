package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

type Environment struct {
	Host string
	Port string
	User string

	KnownHosts string
	PrivateKey string
	Passphrase string

	Command string
}
type Client struct {
	Config     *Environment
	Connection *ssh.Client
	Session    *ssh.Session
}

func getInput(input string) string {
	return os.Getenv("INPUT_" + strings.ToUpper(input))
}

func prepare() {
	err := godotenv.Load()
	er("Error loading .env file: ", err)
}

func config(e *Environment) *ssh.ClientConfig {
	hostKey, err := knownhosts.New(e.KnownHosts)
	er("Failed to load known_hosts", err)

	privateKey, err := os.ReadFile(e.PrivateKey)
	er("Failed to load private key", err)

	signer, err := ssh.ParsePrivateKeyWithPassphrase(privateKey, []byte(e.Passphrase))

	return &ssh.ClientConfig{
		User: e.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: hostKey,
	}
}

func er(message string, err error) {
	if err != nil {
		log.Fatal(message, err.Error())
	}
}

func (e *Environment) newClient() *Client {
	addr := strings.Join([]string{e.Host, e.Port}, ":")

	connection, err := ssh.Dial("tcp", addr, config(e))
	er("Failed to dial: ", err)

	session, err := connection.NewSession()
	er("Failed to create session: ", err)

	return &Client{
		Config:     e,
		Connection: connection,
		Session:    session,
	}
}

func main() {

	prepare()

	env := &Environment{
		Host:       getInput("host"),
		Port:       getInput("port"),
		User:       getInput("user"),
		PrivateKey: getInput("private_key"),
		Passphrase: getInput("passphrase"),
		KnownHosts: getInput("known_hosts"),
		Command:    getInput("command"),
	}

	client := env.newClient()
	defer client.Connection.Close()
	defer client.Session.Close()

	s := client.Session
	var b bytes.Buffer
	s.Stdout = &b

	err := s.Run(env.Command)
	er("Failed to run: ", err)

	fmt.Print(b.String())
}
