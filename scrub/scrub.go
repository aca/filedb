package scrub

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"

	"github.com/aca/farchive/types"
	"github.com/cespare/xxhash"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/xtdlib/filepathx"
	"github.com/xtdlib/sqlitex"
	"github.com/xtdlib/try"
)

type CommandOpt struct {
	FileA string
	FileB string
}

var schema = `
CREATE TABLE IF NOT EXISTS file (
	path TEXT PRIMARY KEY,
	abs TEXT,
	size INTEGER,
	hash TEXT,
	modifiedAt INTEGER,
	validatedAt INTEGER
)
`

type File struct {
	Path string
	// Info os.FileInfo
	Info fs.FileInfo
}

func Command() *cobra.Command {
	f := &CommandOpt{}

	cmd := &cobra.Command{
		Use: "scrub",
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

func readSize(file string) (int64, error) {
	fi, err := os.Stat(file)
	if err != nil {
		return 0, err
	}
	return fi.Size(), nil
}

var db *sqlitex.DB
var dbopt = "?cache=shared&mode=rwc&_busy_timeout=5000&_journal_mode=WAL"

func Run(opt *CommandOpt) error {
	dbname := "farchive.db"

	db = sqlitex.New(dbname + dbopt)
	db.MustExec(schema)

	// ctx := context.Background()

	files := make([]*File, 0, 10000)

	log.Println("walk")

	err := filepathx.WalkDir(".",
		func(path string, d fs.DirEntry, err error) error {
			if path == dbname {
				return nil
			}

			if err != nil {
				return err
			}

			if !filepathx.IsFile(d.Type()) {
				return nil
			}


			switch filepathx.Base(path) {
			case "farchive.db":
				log.Println("skip")
				return nil
			case "snapraid.content":
				log.Println("skip")
				return nil
			}

			info, err := d.Info()
			if err != nil {
				return err
			}

			// HashFile(&File{
			// 	Path: path,
			// 	Info: info,
			// })

			files = append(files, &File{
				Path: path,
				Info: info,
			})

			return nil
		})
	if err != nil {
		log.Fatal(err)
	}

	files = lo.Shuffle(files)

	for _, file := range files {
		log.Println("path: ", file.Path)
		HashFile(file)
	}
	return nil
}

func HashFile(file *File) {
	row := types.Row{}
	path := file.Path
	err := db.Get(&row, "SELECT * FROM file WHERE path = ?", path)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		panic(err)
	}
	hashnew := XXHash(path)
	if row.Hash != hashnew {
		log.Fatal("Hash not match", path, row.Hash, hashnew)
	}
}

func XXHash(path string) string {
	hash := xxhash.New()
	f := try.E1(os.Open(path))
	try.E1(io.Copy(hash, f))
	hashnew := fmt.Sprintf("%016x", hash.Sum64())
	hash.Reset()
	return hashnew
}
