package main

import (
	"log"
	"os"

	"github.com/aca/farchive/diff"
	"github.com/aca/farchive/run"
	"github.com/spf13/cobra"
)

func main() {
	rootcmd, err := newRootCmd(os.Args)
	if err != nil {
		panic(err)
	}
	if err := rootcmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func newRootCmd(args []string) (*cobra.Command, error) {
	versionFlag := false
	cmd := &cobra.Command{
		Use:           os.Args[0],
		SilenceUsage:  false,
		SilenceErrors: false,
	}

	f := cmd.PersistentFlags()
	f.BoolP("verbose", "v", false, "verbose output for debugging purposes")
	f.BoolVar(&versionFlag, "version", false, "print version")
	f.Parse(args)

	cmd.AddCommand(
		run.Command(),
		diff.Command(),
	)

	return cmd, nil
}
