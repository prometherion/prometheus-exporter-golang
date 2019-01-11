package signature

type TaskSignature interface {
	// The signature ID, used to identify the task during unmarshal
	Id() string
	// Serialized bytes to send to the queue persistence layer
	Bytes() (byte []byte, err error)
}

