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

func Run(opt *CommandOpt) error {
	rows := []types.Row{}
	db := sqlitex.New(opt.DB1)
	db2 := sqlitex.New(opt.DB2)

	db.MustSelect(&rows, `select * from file`)

	for _, row := range rows {
		row2 := types.Row{}
		db2.MustGet(&row2, `select * from file where path = ?`, row.Path)

		if row.Hash != row2.Hash {
			log.Fatal("hash not equal", row.Abs, row2.Abs)
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
