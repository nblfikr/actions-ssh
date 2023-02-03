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

func getInput(input string) string {
	return os.Getenv("INPUT_" + strings.ToUpper(input))
}

func prepare() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func shell(cmd string, s *ssh.Session) {
	var b bytes.Buffer
	s.Stdout = &b

	err := s.Run(cmd)
	if err != nil {
		log.Fatal("Failed to run: " + err.Error())
	}

	fmt.Println(b.String())
}

type Environment struct {
	Host string
	Port string
	User string

	KnownHosts string
	PrivateKey string
	Passphrase string

	Command string
}

func config(e *Environment) *ssh.ClientConfig {
	hostKey, err := knownhosts.New(e.KnownHosts)
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := os.ReadFile(e.PrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	signer, err := ssh.ParsePrivateKeyWithPassphrase(privateKey, []byte(e.Passphrase))

	return &ssh.ClientConfig{
		User: e.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: hostKey,
	}
}

func (e *Environment) dial() *ssh.Client {
	addr := strings.Join([]string{e.Host, e.Port}, ":")

	client, err := ssh.Dial("tcp", addr, config(e))
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}

	return client
}

func main() {

	prepare()

	input := Environment{
		Host:       getInput("host"),
		Port:       getInput("port"),
		User:       getInput("user"),
		PrivateKey: getInput("private_key"),
		Passphrase: getInput("passphrase"),
		KnownHosts: getInput("known_hosts"),
		Command:    getInput("command"),
	}

	client := input.dial()
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()

	shell(input.Command, session)
}
