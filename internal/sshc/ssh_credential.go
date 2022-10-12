package sshc

import (
	"errors"
	"os"

	"golang.org/x/crypto/ssh"
)

// ErrEmptyCredential ...
var (
	ErrEmptyCredential = errors.New("no password nor private key provided")
)

// Credential use pwd or private key
type Credential struct {
	Username, Pwd string
	PEMPrivateKey []byte
	PrivateKeyPwd []byte
	conf          *ssh.ClientConfig
}

// NewCredentialWithPwd ...
func NewCredentialWithPwd(user, pwd string) (*Credential, error) {
	c := &Credential{
		Username: user,
		Pwd:      pwd,
	}
	return newCredentia(c)
}

// NewCredentialWithKeyfile ...
func NewCredentialWithKeyfile(user, filename, pwd string) (*Credential, error) {
	privateKey, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return NewCredentialWithKeyData(user, privateKey, pwd)
}

// NewCredentialWithKeyData ...
func NewCredentialWithKeyData(user string, data []byte, pwd string) (*Credential, error) {
	c := &Credential{
		Username:      user,
		PEMPrivateKey: data,
	}
	if pwd != "" {
		c.PrivateKeyPwd = []byte(pwd)
	}
	return newCredentia(c)
}

func newCredentia(c *Credential) (*Credential, error) {
	err := c.init()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Credential) init() error {
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
