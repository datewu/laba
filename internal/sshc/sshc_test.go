package sshc_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/datewu/gtea/utils"
	"github.com/datewu/laba/internal/sshc"
)

func TestAll(t *testing.T) {
	if utils.InGithubCI() {
		return
	}
	fn := filepath.Join(os.Getenv("HOME"), ".ssh/id_rsa")
	cred, err := sshc.NewCredentialWithKeyfile("", fn, "")
	if err != nil {
		t.Fatal(err)
	}
	host := sshc.NewTarget("node1", 36000, *cred)
	err = host.Connect()
	if err != nil {
		t.Fatal(err)
	}
	defer host.Close()
	for i := 0; i < 10; i++ {
		go func(n int) {
			cmd := sshc.Cmd{
				Command: fmt.Sprintf(`echo '%d seesion run: sleep %ds' &&
				 sleep %d && date && echo 'go(%d) exit'`, n, n%5, n%5, n),
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
