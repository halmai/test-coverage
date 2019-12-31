package filereader

import (
	"../types"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

func ReadCommitFiles(repoName string) ([]types.CommitData, error) {
	var ret []types.CommitData

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
		fileName := dirName + file.Name()
		fmt.Println("processing " + fileName + "...")
		content, err := ioutil.ReadFile(fileName)
		if err != nil {
			return ret, err
		}

		commitData, err := parseFile(fileName, string(content))
		if err != nil {
			return ret, err
		}

		ret = append(ret, commitData)
	}

	return ret, nil
}

func parseFile(fileName, content string) (types.CommitData, error) {
	var ret types.CommitData

	lines := strings.Split(content, "\n")
	var err error

	for _, line := range lines {
		if line == "" {
			continue
		}
		suffix := line[strings.Index(line, ": ")+2:]

		if strings.HasPrefix(line, "commit: ") {
			ret.Commit = suffix
			continue
		}

		if strings.HasPrefix(line, "branch: ") {
			ret.BranchName = suffix
			continue
		}

		if strings.HasPrefix(line, "time: ") {
			ret.TimeStamp, err = parseTimeStamp(strings.Trim(suffix, " \t"))
			if err != nil {
				return ret, err
			}
			continue
		}

		if strings.HasPrefix(line, "min") {
			parts := strings.Split(line, " ")
			if len(parts) < 6 {
				return ret, errors.New("wrong min-max-avg line in commit file " + fileName)
			}

			min, err := parsePercentage(parts[1])
			if err != nil {
				return ret, err
			}
			ret.Min = min

			max, err := parsePercentage(parts[3])
			if err != nil {
				return ret, err
			}
			ret.Max = max

			avg, err := parsePercentage(parts[5])
			if err != nil {
				return ret, err
			}
			ret.Avg = avg

			continue
		}
	}

	return ret, nil
}

func parseTimeStamp(str string) (time.Time, error) {
	format := "20060102-150405"
	return time.Parse(format, str)
}

func parsePercentage(str string) (float64, error) {
	value, err := strconv.ParseFloat(str, 32)
	if err != nil {
		return 0, err
	}

	if value<0 || value>100 {
		return 0, errors.New("incorrect percentage value")
	}

	return value, nil
}

