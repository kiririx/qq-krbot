package dao

import (
	"io/fs"
	"io/ioutil"
	"net/url"
	"sync"
	"time"
)

type tag = string
type files = map[int]string

var FileList map[tag]files

func init() {
	FileList = make(map[tag]files)
	go func() {
		for {
			time.Sleep(time.Second * 1)
			files, _ := ioutil.ReadDir("./photo")
			wg := sync.WaitGroup{}
			for _, file := range files {
				wg.Add(1)
				go func(f fs.FileInfo) {
					if f.IsDir() {
						FileList[f.Name()] = func() map[int]string {
							tagFiles := make(map[int]string)
							diskTagFiles, _ := ioutil.ReadDir("./photo/" + f.Name())
							i := 0
							for _, diskTagFile := range diskTagFiles {
								tagFiles[i] = "photo/" + url.QueryEscape(f.Name()) + "/" + diskTagFile.Name()
								i++
							}
							return tagFiles
						}()
					}
					wg.Done()
				}(file)
			}
			wg.Wait()
		}
	}()
}
