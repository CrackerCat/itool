package cmd

import (
	"testing"
)

func TestShellCmd(t *testing.T) {
	ok, err := iFileExsit("/private/var/mobile/Media/iOSSniffer")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(ok)
}
