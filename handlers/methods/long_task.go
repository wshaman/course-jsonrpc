package methods

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type longTaskParams struct {
	A int `json:"a"`
	B int `json:"b"`
}

type task struct {
	ParamA int
	ParamB int
	Status int
	Result int64
}

var tasks = map[string]task{}

const (
	statusInProgress = 0
	statusDone       = 1
)

func Summ(l json.RawMessage) (interface{}, error) {
	ltp := longTaskParams{}
	if err := json.Unmarshal(l, &ltp); err != nil {
		return nil, errors.New("failed to parse longTask params")
	}
	tskID := fmt.Sprintf("task_%d", time.Now().Unix())
	tasks[tskID] = task{
		ParamA: ltp.A,
		ParamB: ltp.B,
		Status: statusInProgress,
		Result: 0,
	}
	go calc(tskID)
	return tskID, nil
}

func calc(taskID string) {
	time.Sleep(15 * time.Second)
	t := tasks[taskID]
	t.Result = int64(t.ParamA) + int64(t.ParamB)
	t.Status = statusDone
	tasks[taskID] = t
}
