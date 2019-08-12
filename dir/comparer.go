package dir

import (
	"fmt"
	"io"
)

// DuplicateResult gives connection between base and target path when duplicate is found
type DuplicateResult struct {
	hash       string
	basePath   string
	targetPath string
}

// FindDuplicate finds duplicates exist whitin targetDir.
func FindDuplicate(baseDir string, targetDir string) []DuplicateResult {
	var result []DuplicateResult
	baseMap := CollectBySingleChannel(baseDir)
	targetMap := CollectBySingleChannel(targetDir)

	for tarHash, tarPath := range targetMap {
		if basePath, exist := baseMap[tarHash]; exist {
			dr := DuplicateResult{tarHash, basePath, tarPath}
			result = append(result, dr)
		}
	}
	return result
}

// WriteDuplicateResult writes the dupres into a writer
func WriteDuplicateResult(dupres []DuplicateResult, w io.Writer) {
	for _, dr := range dupres {
		fmt.Fprintln(w, dr.basePath)
		fmt.Fprintln(w, "> "+dr.hash)
		fmt.Fprintln(w, "~ "+dr.targetPath)
	}
}
