/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/sapslaj/aquapi/internal/service"
	"github.com/spf13/cobra"
)

var (
	add        bool
	remove     bool
	set        bool
	clear      bool
	actionTags []string
)

// tagsCmd represents the tags command
var tagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "do stuff with tags",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		for _, key := range args {
			imageService := service.NewImagesService()
			image, err := imageService.FindById(key)
			if err != nil {
				log.Fatal(err)
			}
			switch true {
			case add:
				if len(actionTags) == 0 {
					log.Fatal("--tag must be specified")
				}
				for _, newTag := range actionTags {
					err = image.AddTag(newTag)
				}
			case remove:
				if len(actionTags) == 0 {
					log.Fatal("--tag must be specified")
				}
				for _, oldTag := range actionTags {
					err = image.RemoveTag(oldTag)
				}
			case set:
				if len(actionTags) == 0 {
					log.Fatal("--tag must be specified")
				}
				err = image.SetTags(actionTags)
			case clear:
				err = image.SetTags([]string{})
			}
			if err != nil {
				log.Fatal(err)
			}
			err = image.Update()
			if err != nil {
				log.Fatal(err)
			}
			tags, err := image.GetTags()
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("%s\ttags: %v", image.ID, tags)
		}
	},
}

func init() {
	rootCmd.AddCommand(tagsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tagsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tagsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	tagsCmd.Flags().BoolVar(&add, "add", false, "Add a tag")
	tagsCmd.Flags().BoolVar(&remove, "remove", false, "remove a tag")
	tagsCmd.Flags().BoolVar(&set, "set", false, "sets tags")
	tagsCmd.Flags().BoolVar(&clear, "clear", false, "clear a tag")
	tagsCmd.Flags().StringSliceVarP(&actionTags, "tag", "t", []string{}, "Single tag or comma separated list of tags to perform action with")
}
