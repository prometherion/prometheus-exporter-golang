package cmd

import (
	"fmt"
	"github.com/sirupsen/logrus"

	"github.com/prometherion/prometheus-exporter-golang/app/queue"
	"github.com/prometherion/prometheus-exporter-golang/app/signature"
	"github.com/spf13/cobra"
)

const (
	// Flags
	name     = "name"
	image    = "image"
	replicas = "replicas"
)

var appCmd = &cobra.Command{
	Use:   "app",
	Short: "Create a Deployment resource",
	Long:  fmt.Sprintf("Create a %s action and push it in the queue in order to let process it by workers", signature.AppCreateName),
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		t, _ := produceCmd.Flags().GetString(action)
		switch t {
		case create:
			_ = cmd.MarkFlagRequired(image)
			_ = cmd.MarkFlagRequired(replicas)
		case update:
			_ = cmd.MarkFlagRequired(image)
			_ = cmd.MarkFlagRequired(replicas)
		}

		_ = cmd.MarkFlagRequired(name)

		return err
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		// Root flags
		a, _ := cmd.Flags().GetString(queueTag)
		cs, _ := cmd.Flags().GetString(redisConString)
		action, _ := cmd.Flags().GetString(action)
		// App flags
		n, _ := cmd.Flags().GetString(name)
		i, _ := cmd.Flags().GetString(image)
		r, _ := cmd.Flags().GetUint32(replicas)

		logrus.Infof("Attempting connection to %s", cs)
		m, err := queue.New(cs, a)
		if err != nil {
			return err
		}
		logrus.Infof("Opening queue at %s", cs)
		m.Open(action)

		var s signature.TaskSignature
		switch action {
		case create:
			logrus.Infof("Task signature is %s", create)
			s = signature.AppCreate{
				Name: n,
				Image: i,
				ReplicaCount: r,
			}
		case read:
			logrus.Infof("Task signature is %s", read)
			s = signature.AppRead{
				Name: n,
			}
		case update:
			logrus.Infof("Task signature is %s", update)
			s = signature.AppUpdate{
				Name: n,
				Image: i,
				ReplicaCount: r,
			}
		case delete:
			logrus.Infof("Task signature is %s", delete)
			s = signature.AppDelete{
				Name: n,
			}
		}

		defer func() {
			if err == nil {
				logrus.Info("Task has been pushed successfully")
			}
		}()
		return m.Push(s)
	},
}

func init() {
	produceCmd.AddCommand(appCmd)

	appCmd.Flags().String(name, "", "The application identifier")
	appCmd.Flags().String(image, "", "The Docker image of the application")
	appCmd.Flags().Uint32(replicas, 0, "Amount of replicas to spin up")
}
