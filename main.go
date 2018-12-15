// Copyright 2018 softlandia@gmail.com
// auto random files renamer

package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/softlandia/xLib"
)

const num int = 5000

func fileTrust(ext, path string, i os.FileInfo) bool {
	if i.IsDir() { //skip dir
		return false
	}
	if filepath.Ext(path) != ext { //skip files with extention not equal extFileName
		return false
	}
	if !xLib.FileExists(path) {
		return false
	}
	return true
}

func arn(n int, path, extFileName string) (int, error) {
	log.Println("start rename")
	log.Printf("max file count: %v\n", n)
	log.Println("file name mask: " + extFileName)

	r := rand.New(rand.NewSource(int64(time.Now().Second()) + int64(time.Now().Minute())))
	i := 0 //index founded files
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if !fileTrust(extFileName, path, info) {
			return nil
		}
		//file found
		i++
		dir := filepath.Dir(path)
		if dir == "." {
			dir = ""
		}
		fileName := info.Name()
		if strings.Contains(fileName, "#") {
			j := strings.Index(fileName, "#") + 1
			fileName = fileName[j:]
		}

		newPath := dir + fmt.Sprintf("%v#%s", r.Int31n(int32(n)), fileName)
		fmt.Printf("%s\n", newPath)
		os.Rename(path, newPath)
		return nil
	})
	return i, err
}
func main() {
	log.Println("program start")
	log.Println("> arn \"x:\\music\" \".mp3\"")
	log.Println(arn(5000, os.Args[1], os.Args[2]))
}
