package queue

import (
	"time"

	"github.com/prometherion/prometheus-exporter-golang/app/signature"
)

type Queue interface {
	Push(signature signature.TaskSignature) (err error)
	StartConsuming(prefetchLimit int, pollDuration time.Duration) (err error)
}
