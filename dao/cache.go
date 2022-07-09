package dao

import (
	"io/fs"
	"io/ioutil"
	"net/url"
	"sync"
	"time"
)

type concurrentMap[K comparable, V any] struct {
	lock *sync.RWMutex
	v    *map[K]V
}

func SyncMap[K comparable, V any]() *concurrentMap[K, V] {
	return &concurrentMap[K, V]{
		lock: &sync.RWMutex{},
		v: func() *map[K]V {
			m := make(map[K]V)
			return &m
		}(),
	}
}

func (c *concurrentMap[K, V]) Get(k K) V {
	c.lock.RLock()
	_v := (*(c.v))[k]
	c.lock.RUnlock()
	return _v
}

func (c *concurrentMap[K, V]) Put(k K, v V) {
	c.lock.Lock()
	(*(c.v))[k] = v
	c.lock.Unlock()
}

func (c *concurrentMap[K, V]) Len() int {
	c.lock.RLock()
	l := len(*c.v)
	c.lock.RUnlock()
	return l
}

type tag = string
type files = *concurrentMap[int, string]

var FileList *concurrentMap[tag, files]

func init() {
	FileList = SyncMap[tag, files]()
	go func() {
		for {
			time.Sleep(time.Second * 10)
			files, _ := ioutil.ReadDir("./photo")
			wg := sync.WaitGroup{}
			for _, file := range files {
				wg.Add(1)
				go func(f fs.FileInfo) {
					if f.IsDir() {
						FileList.Put(f.Name(), func() *concurrentMap[int, string] {
							tagFiles := SyncMap[int, string]()
							diskTagFiles, _ := ioutil.ReadDir("./photo/" + f.Name())
							i := 0
							for _, diskTagFile := range diskTagFiles {
								tagFiles.Put(i, "photo/"+url.QueryEscape(f.Name())+"/"+diskTagFile.Name())
								i++
							}
							return tagFiles
						}())
					}
					wg.Done()
				}(file)
			}
			wg.Wait()
		}
	}()
}
