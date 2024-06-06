package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/cespare/xxhash/v2"
	"github.com/peterbourgon/ff/v4"
	"github.com/samber/lo"
	"github.com/xtdlib/sqlitex"
	"github.com/xtdlib/try"
)

var db *sqlitex.DB

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

type Row struct {
	Path        string `db:"path"`
	Abs         string `db:"abs"`
	Size        int64  `db:"size"`
	Hash        string `db:"hash"`
	ValidatedAt int64  `db:"validatedAt"`
	ModifiedAt    int64  `db:"modifiedAt"`
}

func main() {
	dbname := "farchive.db"
	db = sqlitex.New(dbname)

	db.MustExec(schema)

	rootcmdFlags := ff.NewFlagSet("farchive")
	worker := rootcmdFlags.Int('j', "worker", 1, "parallel worker count")

	// TODO
	_ = worker
	rootcmd := &ff.Command{
		Name:  "farchive",
		Usage: "farchive [FLAGS] SUBCOMMAND ...",
		Flags: rootcmdFlags,
	}
	_ = rootcmd

	// ctx := context.Background()

	type File struct {
		Path string
		Info os.FileInfo
	}

	files := make([]*File, 0, 10000)

	err := filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if path == dbname {
				return nil
			}

			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			files = append(files, &File{
				Path: path,
				Info: info,
			})

			// files = append(files, info)

			return nil
		})
	if err != nil {
		log.Println(err)
	}

	// log.Println(len(files))
	files = lo.Shuffle(files)

	for _, file := range files {
		row := Row{}
		info := file.Info
		path := file.Path
		err = db.Get(&row, "SELECT * FROM file WHERE path = ?", path)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			panic(err)
		}
		if err == nil {
			if row.Size != info.Size() || row.ModifiedAt != info.ModTime().Unix() {
				hashnew := XXHash(path)
				row.ValidatedAt = time.Now().Unix()
				db.MustExec(`UPDATE file SET size = ?, modifiedAt = ?, hash = ?, validatedAt = ? WHERE path = ?`, info.Size(), info.ModTime().Unix(), hashnew, path, row.ValidatedAt)
				if row.Hash != hashnew {
					log.Println("UPDATE HASH", path, row.Hash, hashnew)
				} else {
					log.Println("file not changed: ", path)
				}
			}
		} else {
			row.Path = path
			row.Abs = try.E1(filepath.Abs(path))
			row.Size = info.Size()
			row.ModifiedAt = info.ModTime().Unix()
			row.Hash = XXHash(path)
			row.ValidatedAt = time.Now().Unix()
			log.Println("new file", path)
			db.MustNamedExec(
				`INSERT INTO file (path, abs, size, modifiedAt, hash, validatedAt) VALUES (:path, :abs, :size, :modifiedAt, :hash, :validatedAt)`, map[string]interface{}{
					"path":        row.Path,
					"abs":         row.Abs,
					"size":        row.Size,
					"modifiedAt":  row.ModifiedAt,
					"hash":        row.Hash,
					"validatedAt": row.ValidatedAt,
				},
			)
		}
	}

	// err = rootcmdCmd.ParseAndRun(ctx, os.Args[1:])
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

func XXHash(path string) string {
	hash := xxhash.New()
	f := try.E1(os.Open(path))
	try.E1(io.Copy(hash, f))
	hashnew := fmt.Sprintf("%016x", hash.Sum64())
	hash.Reset()
	return hashnew
}
