package missing

import (
	"fmt"
	"os"

	"github.com/aca/farchive/types"
	"github.com/spf13/cobra"
	"github.com/xtdlib/sqlitex"
)

type CommandOpt struct {
}

func Command() *cobra.Command {
	f := &CommandOpt{}

	cmd := &cobra.Command{
		Use: "missing",
		// Args:          cobra.ExactArgs(2),
		SilenceUsage:  false,
		SilenceErrors: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			// f.FileA = args[0]
			// f.FileB = args[1]
			return Run(f)
		},
	}

	// flags := cmd.Flags()
	// flags.StringVarP(&f.X, "xxx", "x", "default", "Description")

	return cmd
}

var db *sqlitex.DB
var dbopt = "?cache=shared&mode=rwc&_busy_timeout=5000&_journal_mode=WAL"

func Run(opt *CommandOpt) error {
	dbname := "farchive.db"

	db = sqlitex.New(dbname + dbopt)

	rows := []types.Row{}
	err := db.Select(&rows, "SELECT * FROM file")
	if err != nil {
		panic(err)
	}

	for _, row := range rows {
		_, err := os.Stat(row.Path)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}
