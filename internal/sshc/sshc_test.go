package sshc_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/datewu/gtea/utils"
	"github.com/datewu/laba/internal/sshc"
)

func TestAll(t *testing.T) {
	if utils.InGithubCI() {
		return
	}
	privateKey, err := os.ReadFile(os.Getenv("HOME") + "/.ssh/tx-me")
	if err != nil {
		t.Fatal(err)
	}
	cred := &sshc.Credential{
		PEMPrivateKey: privateKey,
	}
	host := &sshc.Target{
		IP:   "9.135.140.60",
		Port: 36000,
	}
	err = host.Connect(cred)
	if err != nil {
		t.Fatal(err)
	}
	defer host.Close()
	for i := 0; i < 10; i++ {
		go func(n int) {
			cmd := sshc.Cmd{
				Command: fmt.Sprintf(`echo '%d seesion run: going to sleep %ds' &&
				 sleep %d && echo 'go(%d) exit'`, n, n%5, n%5, n),
				Out:    os.Stdout,
				Errout: os.Stderr,
			}
			host.Run(cmd)
		}(i)
	}
	cmd := sshc.Cmd{
		Command: "echo sleep && sleep 10 && whoami && echo done",
		Out:     os.Stdout,
		Errout:  os.Stderr,
	}
	err = host.Run(cmd)
	if err != nil {
		t.Fatal(err)
	}
}
