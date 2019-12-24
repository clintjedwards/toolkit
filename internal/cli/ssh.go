package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh"
)

//https://zaiste.net/executing_commands_via_ssh_using_go/

const port string = "22"

func publicKey(path string) ssh.AuthMethod {
	key, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		panic(err)
	}
	return ssh.PublicKeys(signer)
}

func connectSSH(user, host string) (session *ssh.Session, client *ssh.Client, err error) {

	home, err := homedir.Dir()
	if err != nil {
		return nil, nil, err
	}

	sshkeyPath := fmt.Sprintf("%s/%s/%s", home, ".ssh", "id_rsa")

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			publicKey(sshkeyPath),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err = ssh.Dial("tcp", fmt.Sprintf("%s:%s", host, port), config)
	if err != nil {
		log.Fatal(err)
	}

	// Create sesssion
	session, err = client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}

	return session, client, nil
}

func runCommandsOverSSH(user, host string, commands []string) {

	commands = append(commands, "exit")

	session, client, err := connectSSH(user, host)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()
	defer session.Close()

	// StdinPipe for commands
	stdin, err := session.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	// Enable system stdout
	// Comment these if you uncomment to store in variable
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	// Start remote shell
	err = session.Shell()
	if err != nil {
		log.Fatal(err)
	}

	for _, cmd := range commands {
		_, err = fmt.Fprintf(stdin, "%s\n", cmd)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Wait for sess to finish
	err = session.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
