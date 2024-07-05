package lib

import (
	"fmt"
	"io"
	"os"

	"github.com/cespare/xxhash"
	"github.com/xtdlib/try"
)

func XXHash(path string) string {
	hash := xxhash.New()
	f := try.E1(os.Open(path))
	try.E1(io.Copy(hash, f))
	hashnew := fmt.Sprintf("%016x", hash.Sum64())
	hash.Reset()
	return hashnew
}
