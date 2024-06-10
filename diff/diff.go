package diff

import (
	"log"

	"github.com/aca/farchive/types"
	"github.com/spf13/cobra"
	"github.com/xtdlib/sqlitex"
)

type CommandOpt struct {
	DB1 string
	DB2 string
}

var dbopt = "?cache=shared&mode=rwc&_busy_timeout=5000&journal_mode=WAL"

func Run(opt *CommandOpt) error {
	rows := []types.Row{}

	db := sqlitex.New(opt.DB1 + dbopt)
	db2 := sqlitex.New(opt.DB2 + dbopt)

	db.MustSelect(&rows, `select * from file`)

	for _, row := range rows {
		row2 := types.Row{}
		err := db2.Get(&row2, `select * from file where path = ?`, row.Path)
		if sqlitex.IsErrNoRows(err) {
			log.Println("not found", row.Path)
			continue
		}

		if row.Hash != row2.Hash {
			log.Fatal("hash not equal", row.Abs, row2.Abs, row.Hash, row2.Hash)
		}
	}
	return nil
}

func Command() *cobra.Command {
	f := &CommandOpt{}
	cmd := &cobra.Command{
		Use:           "diff",
		Args:          cobra.ExactArgs(2),
		SilenceUsage:  false,
		SilenceErrors: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			f.DB1 = args[0]
			f.DB2 = args[1]
			return Run(f)
		},
	}

	return cmd
}
