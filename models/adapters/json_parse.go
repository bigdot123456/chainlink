package adapters

import (
	"errors"

	simplejson "github.com/bitly/go-simplejson"
	null "gopkg.in/guregu/null.v3"
)

type JsonParse struct {
	Path []string `json:"path"`
}

func (self *JsonParse) Perform(input RunResult) RunResult {
	js, err := simplejson.NewJson([]byte(input.Value()))
	if err != nil {
		return RunResult{Error: err}
	}

	js, err = checkEarlyPath(js, self.Path)
	if err != nil {
		return RunResult{Error: err}
	}

	rval, ok := js.CheckGet(self.Path[len(self.Path)-1])
	if !ok {
		return RunResult{}
	}

	return RunResult{
		Output: map[string]null.String{"value": null.StringFrom(rval.MustString())},
	}
}

func checkEarlyPath(js *simplejson.Json, path []string) (*simplejson.Json, error) {
	var ok bool
	for _, k := range path[:len(path)-1] {
		js, ok = js.CheckGet(k)
		if !ok {
			return js, errors.New("No value could be found for the key '" + k + "'")
		}
	}
	return js, nil
}