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

type HashPath struct {
	hash string
	path string
}

// PrintDirContent get the hash-path values of all files in the directory.
func PrintDirContent(directory string) map[string]string {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}

	result := make(map[string]string)
	//c := make(chan map[string]string)
	cHash := make(chan HashPath)
	//	defer close(cpath)
	//var wg sync.WaitGroup

	ndirectory := 0
	nfiles := 0
	for _, f := range files {
		path := directory + "/" + f.Name()
		if f.IsDir() {
			ndirectory = ndirectory + 1
			//wg.Add(1)
			//go func() {
			submap := PrintDirContent(path)
			//	c <- submap
			//}()
			for kHash, vPath := range submap {
				if val, exist := result[kHash]; exist {
					result[kHash] = val + ", " + vPath
				} else {
					result[kHash] = vPath
				}
			}
			//cpath <- path
		} else {
			nfiles = nfiles + 1
			go func() {
				cHash <- HashPath{getHash(path), path}
			}()
			/*hash := getHash(path)
			if val, exist := result[hash]; exist {
				result[hash] = val + ", " + path
			} else {
				result[hash] = path
			}*/
			// fmt.Println(":: ", path, hash)
		}
	}

	if nfiles == 0 {
		close(cHash)
	}

	nreceived := 0
	for hp := range cHash {
		if val, exist := result[hp.hash]; exist {
			result[hp.hash] = val + ", " + hp.path
		} else {
			result[hp.hash] = hp.path
		}
		nreceived = nreceived + 1
		if nreceived == nfiles {
			close(cHash)
		}
	}
	/*
		if ndirectory == 0 {
			close(cHash)
		}
			for submap := range c {
				nreceived = nreceived + 1
				//fmt.Println("Received something from:", p)

				for kHash, vPath := range submap {
					if val, exist := result[kHash]; exist {
						result[kHash] = val + ", " + vPath
					} else {
						result[kHash] = vPath
					}
				}
				if nreceived == ndirectory {
					close(c)
				}
			}*/

	return result
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
		//fmt.Println(filePath, totalByte)
		hasher.Write(b1)
	}
	return hex.EncodeToString(hasher.Sum(nil))
}
