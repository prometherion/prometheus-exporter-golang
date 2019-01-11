package signature

import (
	"encoding/json"
)

const AppCreateName = "V1_Deployment_Create"

type AppCreate struct {
	Name	 	 string
	Image        string
	ReplicaCount uint32
}

func (AppCreate) Id() string {
	return AppCreateName
}

func (d AppCreate) Bytes() (byte []byte, err error) {
	return json.Marshal(struct {
		AppCreate
		Id string
	}{
		AppCreate: d,
		Id:        d.Id(),
	})
}