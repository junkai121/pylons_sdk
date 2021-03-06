package fixturetest

import (
	testing "github.com/Pylons-tech/pylons_sdk/cmd/fixtures_test/evtesting"
)

// Status is a type to manage work queue status
type Status int

// describes the constant values that describe the status of work queue item
const (
	NotStarted = iota
	InProgress
	Done
)

// QueueItem is a struct to manage work queue item
type QueueItem struct {
	fixtureFileName string
	idx             int
	stepID          string
	status          Status // NotStarted | InProgress | Done
}

var workQueues []QueueItem

// GetQueueID get queue index from scenario file, step index and step id
func GetQueueID(file string, idx int, stepID string) int {
	for i, work := range workQueues {
		if work.fixtureFileName == file && work.stepID == stepID {
			return i
		}
	}
	return -1
}

// GoodToGoForStep check if a step is ready to go
func GoodToGoForStep(file string, idx int, step FixtureStep, t *testing.T) bool {
	for _, condition := range step.RunAfter.PreCondition {
		queID := GetQueueID(file, idx, condition)
		if queID == -1 {
			t.Fatal("No WorkQueue found from specified param ID=", condition, "idx=", idx, "file=", file, workQueues)
		}
		work := workQueues[queID]
		if work.status != Done {
			return false
		}
	}
	return true
}

// UpdateWorkQueueStatus check if a step is ready to go
func UpdateWorkQueueStatus(file string, idx int, fixtureSteps []FixtureStep, targetStatus Status, t *testing.T) {
	step := fixtureSteps[idx]
	queID := GetQueueID(file, idx, step.ID)
	if queID == -1 {
		t.Fatal("No WorkQueue found from specified param ID=", step.ID, "idx=", idx, "file=", file, workQueues)
	}
	switch targetStatus {
	case InProgress:
		if GoodToGoForStep(file, idx, step, t) { // status can move forward only when previous condtions are met
			workQueues[queID].status = InProgress
			ProcessSingleFixtureQueueItem(file, idx, fixtureSteps, t)
		}
	case Done:
		workQueues[queID].status = Done
		if FixtureTestOpts.IsParallel {
			for sidx, sstep := range fixtureSteps {
				squeID := GetQueueID(file, sidx, sstep.ID)
				if workQueues[squeID].status == NotStarted {
					UpdateWorkQueueStatus(file, sidx, fixtureSteps, InProgress, t)
				}
			}
		}
	}
}
