package main

import (
	// "crypto/rand"
	// "encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	// "ggapi/db"
	// "ggapi/model"
	// "reflect"
	//"golang.org/x/crypto/bcrypt"
)

func main() {
	files, _ := filepath.Glob("view/users/*")
	for _, f := range files {
		fmt.Println(os.Remove(f))
	}

}
