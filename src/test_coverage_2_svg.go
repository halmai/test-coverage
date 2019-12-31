package main

import (
	"./filereader"
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"./types"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("This program reads the name of the remote branches from the stdin.")
		fmt.Println("It processes the files residing in the ./<repository-name> directory")
		fmt.Println("From that directory, deletes all the files that contain a branch name which doesn't exist ")
		fmt.Println("and combines the min-max-avg line of the rest of them into a chart.svg file within that directory.")
		fmt.Println("Usage: run test_coverage_2_svg repo_name < git branch -r")
		os.Exit(1)
	}

	repoName := os.Args[1]

	branchNames, err := readBranchNames()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	commitDatas, err := filereader.ReadCommitFiles(repoName)
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}

	_ = branchNames
	// don't delete branches yet. Later on, when SVG becomes crowded, we will delete old data.
	// commitDatas, err = deleteCommitFiles(commitDatas, branchNames)

	fmt.Println("commitDatas", commitDatas)

	svg := createSvg()
	fmt.Println(svg)
}

func createSvg() string {
	ret := `
<?xml version="1.0"?>
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.0//EN" "http://www.w3.org/TR/2001/REC-SVG-20010904/DTD/svg10.dtd">
<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="600" height="300">
	<rect width="300" height="100" style="fill:#fff;stroke:#000099" />
	
	<!-- a text in the same place -->
	<text y="90" style="stroke:#000099;fill:#000099;font-size:24;">chalmai/FS-1234</text>

</svg>

`
	return ret
}

func deleteCommitFiles(commitDatas []types.CommitData, branchNames []string) ([]types.CommitData, error) {
	branches := make(map[string]bool)

	for _, branchName := range branchNames {
		branches[branchName] = true
	}

	// todo: delete file if branch does not exist

	return commitDatas, nil

}

// readBranchNames reads the branch names from stdin,
// requires and removes the "origin/" prefix
// and ignores the HEAD alias
func readBranchNames() ([]string, error) {
	var ret []string
	var emptyList []string

	prefix := "origin/"

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		branchName := strings.Trim(scanner.Text(), " ")

		if !strings.HasPrefix(branchName, prefix) {
			return emptyList, errors.New("branch name misses prefix: " + branchName)
		}

		if strings.HasPrefix(branchName, prefix + "HEAD -> ") {
			continue
		}

		ret = append(ret, strings.Replace(branchName, prefix, "", 1))
	}

	if err := scanner.Err(); err != nil {
		return emptyList, err
	}

	if len(ret) == 0 {
		return emptyList, errors.New("no branches found - probably wrong input")
	}

	return ret, nil
}
