package main

import (
	"albadr/dupfinder/dir"
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	target := "D:/Photos"
	start := time.Now()
	startFormatted := start.Format("2006-01-02 15.04.05")
	resultPath := "Result" + startFormatted + ".dat"

	fmt.Printf("Start hashing %s outputting to %s\n", target, resultPath)
	outFile, err := os.Create(resultPath)
	if err != nil {
		fmt.Println(err)
	}
	defer outFile.Close()

	allmap := dir.CollectBySingleChannel(target)
	w := bufio.NewWriter(outFile)
	for kHash, vPath := range allmap {
		line := kHash + ", " + vPath
		fmt.Fprintln(w, line)
	}
	w.Flush()
	elapsed := time.Since(start)
	log.Printf("Hashing %s took %s\n", target, elapsed)
}
