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
	base := "D:/Photos2"
	target := "L:/F-recovery"
	start := time.Now()
	startFormatted := start.Format("2006-01-02 15.04.05")

	/*
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
	*/

	dupPath := "Dup Result" + startFormatted + ".dat"
	dupResult := dir.FindDuplicate(base, target)

	drOutFile, err := os.Create(dupPath)
	if err != nil {
		fmt.Println(err)
	}
	defer drOutFile.Close()
	drWriter := bufio.NewWriter(drOutFile)
	dir.WriteDuplicateResult(dupResult, drWriter)
	drWriter.Flush()

	dir.MoveTarget(dupResult)

	elapsed := time.Since(start)
	//log.Printf("Hashing %s took %s\n", target, elapsed)
	log.Printf("Hashing %s and %s took %s\n", base, target, elapsed)
}
