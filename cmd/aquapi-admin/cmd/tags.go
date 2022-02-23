/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/sapslaj/aquapi/internal/aquapics"
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
			object, err := aquapics.GetImageObject(key)
			if err != nil {
				log.Fatal(err)
			}
			switch true {
			case add:
				if len(actionTags) == 0 {
					log.Fatal("--tag must be specified")
				}
				for _, newTag := range actionTags {
					err = aquapics.AddTag(object, newTag)
				}
			case remove:
				if len(actionTags) == 0 {
					log.Fatal("--tag must be specified")
				}
				for _, oldTag := range actionTags {
					err = aquapics.RemoveTag(object, oldTag)
				}
			case set:
				if len(actionTags) == 0 {
					log.Fatal("--tag must be specified")
				}
				err = aquapics.SetTags(object, actionTags)
			case clear:
				err = aquapics.SetTags(object, []string{})
			}
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
