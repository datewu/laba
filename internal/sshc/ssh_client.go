package sshc

import (
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/ssh"
)

// ErrEmptyClient ...
var (
	ErrEmptyClient = errors.New("no client connect")
)

// Target reprenet a ssh server
type Target struct {
	Addr   string // IP or hostname
	Port   int
	client *ssh.Client
}

// Connect to target, should call Close afterwards
func (t *Target) Connect(c *Credential) error {
	if t.Port == 0 {
		t.Port = 22
	}
	if c.conf == nil {
		err := c.init()
		if err != nil {
			return err
		}
	}
	addr := fmt.Sprintf("%s:%d", t.Addr, t.Port)
	client, err := ssh.Dial("tcp", addr, c.conf)
	if err != nil {
		return err
	}
	t.client = client
	return nil

}

// Close connection
func (t *Target) Close() error {
	return t.client.Close()
}

type Cmd struct {
	Command     string
	Out, Errout io.Writer
}

// Run commands
func (t *Target) Run(c Cmd) error {
	if t.client == nil {
		return ErrEmptyClient
	}
	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := t.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.
	session.Stdout = c.Out
	session.Stderr = c.Errout
	if err := session.Run(c.Command); err != nil {
		return err
	}
	// return session.Wait()
	return nil
}
