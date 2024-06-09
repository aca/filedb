package main

import (
	"log"
	"os"

	"github.com/xtdlib/try"
)

func main() {
	err := X()
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("ok")
		}
		log.Println(err)
	}
}

func X() (ferr error) {
	defer try.Handle(&ferr)
	f := try.E1(os.Open("notexist"))
	log.Println(f.Name())
	return nil
}
