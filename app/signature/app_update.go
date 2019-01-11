package signature

import "encoding/json"

const AppUpdateName = "V1_Deployment_Update"

type AppUpdate struct {
	Name		 string
	Image        string
	ReplicaCount uint32
}

func (AppUpdate) Id() string {
	return AppUpdateName
}

func (d AppUpdate) Bytes() (byte []byte, err error) {
	return json.Marshal(struct {
		AppUpdate
		Id string
	}{
		AppUpdate: d,
		Id:      d.Id(),
	})
}