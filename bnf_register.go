package eval

import (
	"sync"
	"encoding/json"
	"fmt"
)

type RegisterDataItem struct {
	Index   string `json:"index"`
	Content string `json:"exp"`
	Value   *Value
}

func (rdi *RegisterDataItem) UnmarshalJSON(data []byte) error {
	m := make(map[string]string)
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	rdi.Index = m["index"]
	rdi.Content = m["exp"]
	if rdi.Content == "" {
		return fmt.Errorf("empty express %+v", m)
	}
	var val Value
	err = parser.ParseString(rdi.Content, &val)
	if err != nil {
		return err
	}
	rdi.Value = &val
	return nil
}



var evalMapLocker sync.RWMutex
var evalMap map[string]RegisterDataItem

func RegisterEvalExpress(data []byte) error {
	m := make(map[string]RegisterDataItem)
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	evalMapLocker.Lock()
	defer evalMapLocker.Unlock()
	evalMap = m
	return nil
}

func findExp(expKey string) (val Value, bool bool) {
	evalMapLocker.RLock()
	defer evalMapLocker.RUnlock()

	item, exist := evalMap[expKey]
	if !exist {
		return Value{}, false
	}
	return *item.Value, true
}

func findExpMust(expKey string) Value {
	v, exist := findExp(expKey)
	if !exist {
		panicFor(fmt.Sprintf("not support exp key: %s", expKey))
	}
	return v
}