package signature

import "encoding/json"

const AppReadName = "V1_Deployment_Update"

type AppRead struct {
	Name string
}

func (AppRead) Id() string {
	return AppReadName
}

func (d AppRead) Bytes() (byte []byte, err error) {
	return json.Marshal(struct {
		AppRead
		Name string
	}{
		Name:      d.Id(),
	})
}