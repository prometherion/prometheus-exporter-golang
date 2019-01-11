package queue

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/adjust/rmq"
	"github.com/prometherion/prometheus-exporter-golang/app/signature"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

const (
	namespace = "prometherion"
	subSystem = "worker"
)

func NewConsumer(queue string) Consumer {
	return Consumer{
		unmarshalMetric: *promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subSystem,
			Name:      "tasks_error_unmarshal",
			Help:      "Total number of rejected tasks due to unmarshal error",
			ConstLabels: map[string]string{
				"queue": queue,
			},
		}, nil),
		rejectedMetric: *promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subSystem,
			Name:      "tasks_error_rejected",
			Help:      "Total number of unrecognized tasks",
			ConstLabels: map[string]string{
				"queue": queue,
			},
		}, []string{"name"}),
		taskMetric: *promauto.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: subSystem,
			Name:      "tasks_completed",
			Help:      "Histogram of tasks completed",
			Buckets: []float64{
				float64(1),
				float64(5),
				float64(10),
				float64(15),
				float64(30),
				float64(60),
				float64(120),
			},
			ConstLabels: map[string]string{
				"queue": queue,
			},
		}, []string{"task", "success"}),

	}
}

type Consumer struct {
	unmarshalMetric prometheus.CounterVec
	rejectedMetric  prometheus.CounterVec
	taskMetric      prometheus.HistogramVec
}

func (c Consumer) Consume(delivery rmq.Delivery) {
	logrus.Info("Consuming a new delivery")
	startTime := time.Now()

	var err error
	var task interface{}

	if err := json.Unmarshal([]byte(delivery.Payload()), &task); err != nil {
		c.unmarshalMetric.WithLabelValues().Inc()
		logrus.Warningf("Rejecting task due to error (%s)", err.Error())
		delivery.Reject()
		return
	}

	m := task.(map[string]interface{})
	id := m["Id"].(string)

	logrus.Infof("Executing delivery unmarshal of task %s", id)

	switch id {
	case signature.AppCreateName:
		var s signature.AppCreate
		_ = json.Unmarshal([]byte(delivery.Payload()), &s)
		c.consume(s)
	case signature.AppUpdateName:
		var s signature.AppUpdate
		_ = json.Unmarshal([]byte(delivery.Payload()), &s)
		c.consume(s)
	case signature.AppDeleteName:
		var s signature.AppDelete
		_ = json.Unmarshal([]byte(delivery.Payload()), &s)
		c.consume(s)
	// You will notice that signature.AppReadName is missing:
	// this is just to ensure that we'll have some unrecognized signature metrics to populate
	default:
		logrus.Warningf("Rejecting task due to unrecognized signature (%s)", id)
		c.rejectedMetric.WithLabelValues(id).Inc()
		delivery.Reject()
		return
	}
	delivery.Ack()

	// Just faking a random boolean in order to fake a error during task processing
	if rand.Intn(2) == 0 {
		err = fmt.Errorf("you unlucky guy")
		logrus.Errorf(err.Error())
	}

	duration := time.Now().Sub(startTime)
	c.taskMetric.WithLabelValues(id, fmt.Sprintf("%t", err == nil)).Observe(duration.Seconds())

	logrus.Warningf("Task %s performed in %.4f seconds", id, duration.Seconds())
}

func (c Consumer) consume(signature signature.TaskSignature) {
	b, _ := signature.Bytes()
	logrus.Info("Faking processing, waiting some seconds")
	w := rand.Intn(len(b))
	time.Sleep(time.Duration(w) * time.Second)
}