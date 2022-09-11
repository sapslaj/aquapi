package cmd

import (
	"log"

	"github.com/sapslaj/aquapi/internal/service"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get a thing",
	Long:  `todo`,
	Run: func(cmd *cobra.Command, args []string) {
		imageService := service.NewImagesService()
		if len(args) == 0 {
			image, err := imageService.GetRandomImageFilterTags(nil, nil)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("%s\ttags: %v", image.ID, image.Tags)
		} else {
			for _, key := range args {
				image, err := imageService.FindById(key)
				if err != nil {
					log.Fatal(err)
				}
				log.Printf("%s\ttags: %v", image.ID, image.Tags)
			}
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
