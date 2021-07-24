package methods

import (
	"encoding/json"
	"errors"
)

type longTaskStatusParams struct {
	ID string `json:"id"`
}

type longTaskStatusResponse struct {
	A    int   `json:"a"`
	B    int   `json:"b"`
	Summ int64 `json:"summ"`
}

func SummStatus(l json.RawMessage) (interface{}, error) {
	ltp := longTaskStatusParams{}
	if err := json.Unmarshal(l, &ltp); err != nil {
		return nil, errors.New("failed to parse longStatusTask params")
	}
	v, ok := tasks[ltp.ID]
	if !ok {
		return nil, errors.New("task not found")
	}
	if v.Status == statusInProgress {
		return "operation in progress", nil
	}

	res := longTaskStatusResponse{
		A:    v.ParamA,
		B:    v.ParamB,
		Summ: v.Result,
	}
	delete(tasks, ltp.ID)
	return res, nil
}
