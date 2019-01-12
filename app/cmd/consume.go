package cmd

import (
	"net/http"
	"os"
	"time"

	"github.com/prometherion/prometheus-exporter-golang/app/queue"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var consumeCmd = &cobra.Command{
	Use:   "consume",
	Short: "Start consuming tasks of a selected queue",
	Long:  "Spawns multiple workers as goroutines to process tasks belonging to the selected queue.",
	Run: func(cmd *cobra.Command, args []string) {
		// Root flags
		t, _ := cmd.Flags().GetString(queueTag)
		cs, _ := cmd.Flags().GetString(redisConString)
		task, _ := cmd.Flags().GetString(action)

		m, err := queue.New(cs, t)
		if err != nil {
			logrus.Error(err.Error())
			os.Exit(1)
		}
		logrus.Infof("Consuming tasks of queue %s", task)
		m.Open(task)

		err = m.StartConsuming(10, time.Second)
		if err != nil {
			logrus.Error(err.Error())
			os.Exit(1)
		}

		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":9000", nil); err != nil {
			logrus.Errorf("Cannot start prometheus metrics endpoint (%s)", err.Error())
		}
		select {}
	},
}

func init() {
	rootCmd.AddCommand(consumeCmd)

	consumeCmd.PersistentFlags().String(action, "", "Select the CRUD operation to consume")
	_ = consumeCmd.MarkPersistentFlagRequired(action)
}
