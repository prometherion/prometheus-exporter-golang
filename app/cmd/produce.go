package cmd

import (
	"github.com/spf13/cobra"
)

var produceCmd = &cobra.Command{
	Use:   "produce",
	Short: "Produce a action",
	Long:  "Produce an available action in order to be consumed by worker",
}

func init() {
	rootCmd.AddCommand(produceCmd)

	produceCmd.PersistentFlags().String(action, "", "Select the CRUD operation to produce")
	_ = produceCmd.MarkPersistentFlagRequired(action)
}
