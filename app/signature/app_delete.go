package signature

import "encoding/json"

const AppDeleteName = "V1_Deployment_Delete"

type AppDelete struct {
	Name 		 string
}

func (AppDelete) Id() string {
	return AppDeleteName
}

func (d AppDelete) Bytes() (byte []byte, err error) {
	return json.Marshal(struct {
		AppDelete
		Id string
	}{
		AppDelete: d,
		Id:        d.Id(),
	})
}