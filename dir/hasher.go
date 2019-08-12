package dir

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

const MAX_BYTES_TO_HASH = 10000

// HashPath is a tupple of hash and path string
type HashPath struct {
	hash string
	path string
}

// HashPaths is a tupple of hash and some paths
type HashPaths struct {
	hash  string
	paths []string
}

// CollectHashOf get the hash-path values of all files in the directory.
func CollectHashOf(directory string) map[string]HashPaths {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}

	result := make(map[string]HashPaths)
	cHash := make(chan HashPath)
	nfiles := 0
	for _, f := range files {
		path := directory + "/" + f.Name()
		if f.IsDir() {
			submap := CollectHashOf(path)
			for kHash, vHp := range submap {
				if val, exist := result[kHash]; exist {
					val.paths = append(val.paths, vHp.paths...)
					result[kHash] = val
				} else {
					paths := []string{path}
					result[kHash] = HashPaths{kHash, paths}
				}
			}
		} else {
			nfiles = nfiles + 1
			go func() {
				cHash <- HashPath{getHash(path), path}
			}()
		}
	}

	if nfiles == 0 {
		close(cHash)
	}

	nreceived := 0
	for hp := range cHash {
		if val, exist := result[hp.hash]; exist {
			val.paths = append(val.paths, hp.path)
			result[hp.hash] = val
		} else {
			result[hp.hash] = HashPaths{hp.hash, []string{hp.path}}
		}
		nreceived = nreceived + 1
		if nreceived == nfiles {
			close(cHash)
		}
	}
	return result
}

//CollectBySingleChannel collects hash using single channel in the implementation
func CollectBySingleChannel(directory string) map[string]HashPaths {
	result := make(map[string]HashPaths)

	jobsChannel := make(chan string, 100)
	hashesChannel := make(chan HashPath, 100)
	paths := collectFilePaths(directory)

	for w := 1; w <= 5; w++ {
		go hashingWorker(jobsChannel, hashesChannel)
	}

	go func() {
		for _, p := range paths {
			jobsChannel <- p
			//fmt.Println(p)
		}
		close(jobsChannel)
		fmt.Println("jobs channel close")
	}()

	for i := 0; i < len(paths); i++ {
		hp := <-hashesChannel
		if val, exist := result[hp.hash]; exist {
			val.paths = append(val.paths, hp.path)
			result[hp.hash] = val
		} else {
			result[hp.hash] = HashPaths{hp.hash, []string{hp.path}}
		}
	}
	return result
}

func collectFilePaths(dir string) []string {
	var result []string
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		path := dir + "/" + f.Name()
		if f.IsDir() {
			result = append(result, collectFilePaths(path)...)
		} else {
			result = append(result, path)
		}
	}
	return result
}

func channeledHash(filePath string, resultChannel chan<- HashPath) {
	hash := getHash(filePath)
	resultChannel <- HashPath{hash, filePath}
}

func hashingWorker(filePaths <-chan string, hashResults chan<- HashPath) {
	for filePath := range filePaths {
		hashResults <- HashPath{getHash(filePath), filePath}
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getHash(filePath string) string {
	f, err := os.Open(filePath)
	check(err)

	hasher := md5.New()

	const bufferSize = 100
	totalByte := 0
	b1 := make([]byte, 5)
	for {
		nRead, err := f.Read(b1)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
		totalByte = totalByte + nRead
		if totalByte > MAX_BYTES_TO_HASH {
			break
		}
		hasher.Write(b1)
	}
	return hex.EncodeToString(hasher.Sum(nil))
}
