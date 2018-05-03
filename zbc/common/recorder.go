package zbcommon

import (
	"encoding/json"
	"io/ioutil"
)

type EventRecorder struct {
	Records []interface{}
}

func (er *EventRecorder) Add(record interface{}) {
	er.Records = append(er.Records, record)
}

func (er *EventRecorder) Dump(filepath string) {
	b, _ := json.Marshal(er)
	ioutil.WriteFile(filepath, b, 0644)
}
