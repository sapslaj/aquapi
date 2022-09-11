package cmd

import (
	"context"

	"github.com/sapslaj/aquapi/internal/maintenance"
	"github.com/spf13/cobra"
)

var syncTableCmd = &cobra.Command{
	Use: "sync-table",
	Run: func(cmd *cobra.Command, args []string) {
		maintenance.SyncImagesBucketToTable(context.TODO())
	},
}

func init() {
	rootCmd.AddCommand(syncTableCmd)
}
