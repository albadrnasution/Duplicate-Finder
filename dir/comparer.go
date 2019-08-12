package dir

import (
	"fmt"
	"io"
)

// DuplicateResult gives connection between base and target path when duplicate is found
type DuplicateResult struct {
	hash       string
	basePath   []string
	targetPath []string
}

// FindDuplicate finds duplicates exist whitin targetDir.
func FindDuplicate(baseDir string, targetDir string) []DuplicateResult {
	var result []DuplicateResult
	baseMap := CollectBySingleChannel(baseDir)
	targetMap := CollectBySingleChannel(targetDir)

	for tarHash, tarHp := range targetMap {
		if baseHp, exist := baseMap[tarHash]; exist {
			dr := DuplicateResult{tarHash, baseHp.paths, tarHp.paths}
			result = append(result, dr)
		}
	}
	return result
}

// WriteDuplicateResult writes the dupres into a writer
func WriteDuplicateResult(dupres []DuplicateResult, w io.Writer) {
	for _, dr := range dupres {
		fmt.Fprintln(w, dr.hash)
		for i, base := range dr.basePath {
			fmt.Fprintln(w, i, "> ", base)
		}
		for i, target := range dr.targetPath {
			fmt.Fprintln(w, i, "+ ", target)
		}
	}
}
