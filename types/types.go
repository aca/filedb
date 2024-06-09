package types

type Row struct {
	Path        string `db:"path"`
	Abs         string `db:"abs"`
	Size        int64  `db:"size"`
	Hash        string `db:"hash"`
	ValidatedAt int64  `db:"validatedAt"`
	ModifiedAt  int64  `db:"modifiedAt"`
}
