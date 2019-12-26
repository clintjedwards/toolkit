package sshutil

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh"
)

//https://zaiste.net/executing_commands_via_ssh_using_go/

const port string = "22"

func getKeyAuth(keyfile string) ssh.AuthMethod {
	key, err := ioutil.ReadFile(keyfile)
	if err != nil {
		panic(err)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		panic(err)
	}
	return ssh.PublicKeys(signer)
}

// connect to specified ssh server; attempts to use ssh-key auth by looking for
// default $HOME/.ssh/id_rsa keyfile location
// Caller should remember to close client and session
func connect(user, host string) (session *ssh.Session, client *ssh.Client, err error) {

	home, err := homedir.Dir()
	if err != nil {
		return nil, nil, err
	}

	sshKeyPathFmt := "%s/%s/%s"
	sshKeyPath := fmt.Sprintf(sshKeyPathFmt, home, ".ssh", "id_rsa")

	_, err = os.Stat(sshKeyPath)
	if os.IsNotExist(err) {
		return nil, nil, fmt.Errorf("could not get ssh key from: %s", sshKeyPath)
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			getKeyAuth(sshKeyPath),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err = ssh.Dial("tcp", fmt.Sprintf("%s:%s", host, port), config)
	if err != nil {
		return nil, nil, fmt.Errorf("could not connect to server: %w", err)
	}

	// A session helps us run multiple commands in a single connection
	session, err = client.NewSession()
	if err != nil {
		return nil, nil, fmt.Errorf("could not establish session: %w", err)
	}

	return session, client, nil
}

// RunCommandsOverSSH establishes a connection with server provided and inputs commands
// in a single session
func RunCommandsOverSSH(hostname string, commands []string) error {

	// We need some way to tell the server we no longer want to send commands
	// and exit the session
	commands = append(commands, "exit")

	hostParts := strings.Split(hostname, "@")

	session, client, err := connect(hostParts[0], hostParts[1])
	if err != nil {
		return fmt.Errorf("could not create connection: %w", err)
	}

	defer client.Close()
	defer session.Close()

	// StdinPipe for commands
	stdin, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("could not get stdin pipe: %w", err)
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	// Start remote shell
	err = session.Shell()
	if err != nil {
		return fmt.Errorf("could not start remote shell: %w", err)
	}

	for _, cmd := range commands {
		_, err = fmt.Fprintf(stdin, "%s\n", cmd)
		if err != nil {
			return fmt.Errorf("could not run command '%s': %w", cmd, err)
		}
	}

	// Wait for sess to finish
	err = session.Wait()
	if err != nil {
		return fmt.Errorf("error in session wait: %w", err)
	}

	return nil
}
