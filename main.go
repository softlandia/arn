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

	"github.com/softlandia/xlib"
)

const num int32 = 50

func fileTrust(ext, path string, i os.FileInfo) bool {
	if i.IsDir() { //skip dir
		return false
	}
	if filepath.Ext(path) != ext { //skip files with extention not equal extFileName
		return false
	}
	if !xlib.FileExists(path) {
		return false
	}
	return true
}

func findFiles(fileList *[]string, path, extFileName string) (int, error) {
	log.Println("start search")
	log.Println("start path: " + path)
	log.Println("file name mask: " + extFileName)

	i := 0 //index founded files
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if !fileTrust(extFileName, path, info) {
			//log.Println(">path: " + path)
			return nil
		}
		//file found
		i++
		*fileList = append(*fileList, path)
		return nil
	})
	return i, err
}

func arn(fileList *[]string, extFileName string) (int, error) {
	log.Println("start rename")
	r := rand.New(rand.NewSource(int64(time.Now().Second()) + int64(time.Now().Minute())))
	//	i := 0 //index founded files
	for _, fn := range *fileList {
		dir := filepath.Dir(fn)
		//log.Println(fn)
		//log.Println(dir)
		newPath := ""
		if strings.Contains(fn, "#") {
			j := strings.Index(fn, "#") + 1
			fn = fn[j:]
		}
		if dir == "." {
			newPath = fmt.Sprintf("%v#%s", r.Int31n(num), filepath.Base(fn))
		} else {
			newPath = fmt.Sprintf("%s\\%v#%s", dir, r.Int31n(num), filepath.Base(fn))
		}
		fmt.Printf("%s\n", newPath)
		os.Rename(fn, newPath)
	}
	return 0, nil
}

func main() {
	log.Println("program start")
	if len(os.Args) == 1 {
		log.Println("> arn \"x:\\music" \".mp3\"")
	}
	if os.Args[1] == "+" {
		//тестовый режим, создаём в каталоге os.Args[2] 2400 файлов
		makeFiles(os.Args[2])
	} else {
		fileList := make([]string, 0, 10)
		fmt.Printf("path to search: '%s'\n", os.Args[1])
		fmt.Printf("extention: '%s'\n", os.Args[2])
		n, err := xlib.FindFilesExt(&fileList, os.Args[1], os.Args[2])
		if err != nil {
			log.Panic(err)
		}
		//i, _ := findFiles(&fileList, os.Args[1], os.Args[2])
		fmt.Printf("founded :%v files\n", n)
		log.Println(arn(&fileList, os.Args[2]))
	}
}

//create foo files
func makeFiles(path string) {
	oFileName := ""
	for i := 0; i < 2400; i++ {
		oFileName = fmt.Sprintf(path+"%v.txt", i)
		fmt.Println("make file: '", oFileName)
		oFile, err := os.Create(oFileName) //Open file to WRITE
		if err != nil {
			fmt.Println("file: " + oFileName + " can't open to write")
		}
		defer oFile.Close()
		fmt.Fprintf(oFile, "%v", i)
	}
}
