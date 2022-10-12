package sshc

import (
	"testing"

	"github.com/datewu/gtea/utils"
)

func TestClient(t *testing.T) {
	if utils.InGithubCI() {
		return
	}
}
