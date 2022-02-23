package cmd

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/sapslaj/aquapi/internal/aquapics"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get a thing",
	Long:  `todo`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, key := range args {
			object, err := aquapics.GetImageObject(key)
			if err != nil {
				log.Fatal(err)
			}
			tags, err := aquapics.GetTags(object)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("%s\ttags: %v", aws.ToString(object.Key), tags)
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
