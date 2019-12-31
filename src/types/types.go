package types

import (
	"fmt"
	"time"
)

type CommitData struct {
	Commit string
	BranchName string
	TimeStamp time.Time
	Min, Max, Avg float64
}

func (cd CommitData) String() string {
	tpl := "Commit: %s, branchName: %s, time: %v, min-max-avg: %f-%f-%f"
	return fmt.Sprintf(tpl, cd.Commit, cd.BranchName, cd.TimeStamp, cd.Min, cd.Max, cd.Avg)
}

