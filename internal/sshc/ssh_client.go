package sshc

import (
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/ssh"
)

// ErrEmptyCredential ...
var ErrEmptyCredential = errors.New("no password nor private key provided")

// Target reprenet a ssh server
type Target struct {
	IP     string
	Port   int
	client *ssh.Client
}

// Connect to target, should call Close afterwards
func (t *Target) Connect(c *Credential) error {
	if t.Port == 0 {
		t.Port = 22
	}
	if c.conf == nil {
		err := c.Init()
		if err != nil {
			return err
		}
	}
	addr := fmt.Sprintf("%s:%d", t.IP, t.Port)
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
	// // Each ClientConn can support multiple interactive sessions,
	// // represented by a Session.
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

// Credential use pwd or private key
type Credential struct {
	Username, Pwd string
	PEMPrivateKey []byte
	PrivateKeyPwd []byte
	conf          *ssh.ClientConfig
}

func (c *Credential) Init() error {
	// var hostKey ssh.PublicKey
	if c.Pwd == "" && c.PEMPrivateKey == nil {
		return ErrEmptyCredential
	}
	if c.Username == "" {
		c.Username = "root"
	}
	conf := &ssh.ClientConfig{
		User: c.Username,
		// HostKeyCallback: ssh.FixedHostKey(hostKey),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	if c.Pwd != "" {
		auth := ssh.Password(c.Pwd)
		conf.Auth = append(conf.Auth, auth)
		c.conf = conf
		return nil // use pwd, ignore Privatekey
	}
	if c.PEMPrivateKey != nil {
		var signer ssh.Signer
		var err error
		if c.PrivateKeyPwd == nil {
			signer, err = ssh.ParsePrivateKey(c.PEMPrivateKey)
		} else {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(
				c.PEMPrivateKey, c.PrivateKeyPwd)
		}
		if err != nil {
			return err
		}
		auth := ssh.PublicKeys(signer)
		conf.Auth = append(conf.Auth, auth)
	}
	c.conf = conf
	return nil
}
