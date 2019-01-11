package signature

import "encoding/json"

const AppReadName = "V1_Deployment_Read"

type AppRead struct {
	Name string
}

func (AppRead) Id() string {
	return AppReadName
}

func (d AppRead) Bytes() (byte []byte, err error) {
	return json.Marshal(struct {
		AppRead
		Id string
	}{
		Id: d.Id(),
	})
}