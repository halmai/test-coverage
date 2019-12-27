package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
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

	commitDatas, err := readCommitFiles(repoName)
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}

	commitDatas, err = deleteCommitFiles(commitDatas, branchNames)

	fmt.Println("commitDatas", commitDatas)
}

func deleteCommitFiles(commitDatas []commitData, branchNames []string) ([]commitData, error) {
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

func readCommitFiles(repoName string) ([]commitData, error) {
	var ret []commitData

	dirName := "./" + repoName + "/"
	f, err := os.Open(dirName)
	if err != nil {
		return ret, err
	}

	files, err := f.Readdir(-1)
	defer f.Close()
	if err != nil {
		return ret, err
	}

	for _, file := range files {
		fmt.Println("processing " + dirName + file.Name() + "...")
		content, err := ioutil.ReadFile(dirName + file.Name())
		if err != nil {
			return ret, err
		}

		commitData, err := parseFile(string(content))
		if err != nil {
			return ret, err
		}

		ret = append(ret, commitData)
	}

	return ret, nil
}

type commitData struct {
	commit string
	branchName string
	timeStamp string
	min, max, avg float64
}

func parseFile(content string) (commitData, error) {
	var ret commitData

	lines := strings.Split(content, "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}
		suffix := line[strings.Index(line, ": "):]

		if strings.HasPrefix(line, "commit: ") {
			ret.commit = suffix
			continue
		}

		if strings.HasPrefix(line, "branch: ") {
			ret.branchName = suffix
			continue
		}

		if strings.HasPrefix(line, "time: ") {
			ret.timeStamp = suffix
			continue
		}

		if strings.HasPrefix(line, "min") {
			parts := strings.Split(line, " ")
			if len(parts) < 6 {
				return ret, errors.New("wrong min-max-avg line in commit file ")
			}

			min, err := strconv.ParseFloat(parts[1], 32)
			if err != nil {
				return ret, err
			}
			ret.min = min

			max, err := strconv.ParseFloat(parts[3], 32)
			if err != nil {
				return ret, err
			}
			ret.max = max

			avg, err := strconv.ParseFloat(parts[5], 32)
			if err != nil {
				return ret, err
			}
			ret.avg = avg

			continue
		}
	}

	return ret, nil
}