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
	//reader := bufio.NewReader(os.Stdin)
	//fmt.Print("Enter text: ")
	//text, _ := reader.ReadString('\n')
	//fmt.Println("Printing directory contents of: ", text)
	//dir.PrintDirContent("d:/Documents/Pribadi/TRAVEL/")
	//args := os.Args[1:]

	target := "D:/Photos"
	start := time.Now()
	fmt.Println("Running ", target, start.Format("2006-01-02 15.04.05"))
	//fmt.Println("Reading", args[1])

	outFile, err := os.Create(start.Format("2006-01-02 15.04.05") + ".txt")
	if err != nil {
		fmt.Println(err)
	}
	defer outFile.Close()

	//target := "D:/Test"
	allmap := dir.PrintDirContent(target)

	w := bufio.NewWriter(outFile)
	for kHash, vPath := range allmap {
		//fmt.Println(kHash, vPath)
		line := kHash + ", " + vPath
		fmt.Fprintln(w, line)
	}
	w.Flush()
	elapsed := time.Since(start)
	log.Printf("Hashing %s took %s", target, elapsed)
}
