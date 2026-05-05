package cli

import (
	"fmt"
	"os"
	"strings"
)

func init() {
	registerCommand("redact", func(args []string, st storeIface) error {
		if len(args) < 3 {
			return fmt.Errorf("usage: envchain redact <project> <passphrase> <text>")
		}
		project := args[0]
		passphrase := args[1]
		text := strings.Join(args[2:], " ")
		return CmdRedact(st.(*storeAdapter).s, project, passphrase, text, os.Stdout)
	})

	registerCommand("redact-file", func(args []string, st storeIface) error {
		if len(args) < 3 {
			return fmt.Errorf("usage: envchain redact-file <project> <passphrase> <file|->")
		}
		project := args[0]
		passphrase := args[1]
		path := args[2]
		return CmdRedactFile(st.(*storeAdapter).s, project, passphrase, path, os.Stdout)
	})
}
