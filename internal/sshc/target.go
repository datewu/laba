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
	Addr string // IP or hostname
	Port int
	Credential
	client *ssh.Client
}

// NewTarget ...
func NewTarget(addr string, port int, c Credential) *Target {
	t := &Target{Addr: addr, Port: port}
	t.Credential = c
	return t
}

// Connect to target, should call Close afterwards
func (t *Target) Connect() error {
	if t.Port == 0 {
		t.Port = 22
	}
	if t.Credential.conf == nil {
		return ErrEmptyCredential
	}
	addr := fmt.Sprintf("%s:%d", t.Addr, t.Port)
	client, err := ssh.Dial("tcp", addr, t.Credential.conf)
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

// Run each clientConn can support multiple interactive sessions.
func (t *Target) Run(c Cmd) error {
	if t.client == nil {
		return ErrEmptyClient
	}
	session, err := t.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	session.Stdout = c.Out
	session.Stderr = c.Errout
	if err := session.Run(c.Command); err != nil {
		return err
	}
	// return session.Wait()
	return nil
}
