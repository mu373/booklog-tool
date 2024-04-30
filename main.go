package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

var (
	location string
	readFile string
	newTag   string
	idType   string
	config   *Config
)

var rootCmd = &cobra.Command{
	Use:   "booklog-tool",
	Short: "A CLI tool to update book information on booklog.jp",
}

var updateLocationCmd = &cobra.Command{
	Use:   "update-location",
	Short: "Update the location for items",
	Run: func(cmd *cobra.Command, args []string) {
		// Read IDs from the file
		file, err := os.Open(readFile)
		if err != nil {
			fmt.Println("Error opening ID file:", err)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			id := scanner.Text()
			var err error

			if idType == "itemid" {
				err = UpdateItemLocationByItemId(id, "loc_"+location)
			} else if idType == "isbn" {
				err = UpdateItemLocationByIsbn(id, "loc_"+location)
			} else {
				fmt.Println("Invalid ID type")
			}

			if err != nil {
				fmt.Printf("Error updating location for ID %s: %v\n", id, err)
			}
			time.Sleep(500 * time.Millisecond)
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading ID file:", err)
		}
	},
}

var addTagCmd = &cobra.Command{
	Use:   "add-tag",
	Short: "Add a tag to items",
	Run: func(cmd *cobra.Command, args []string) {
		// Read IDs from the file
		file, err := os.Open(readFile)
		if err != nil {
			fmt.Println("Error opening ID file:", err)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			id := scanner.Text()
			var err error

			if idType == "itemid" {
				err = AddTagToItemByItemId(id, newTag)
			} else if idType == "isbn" {
				err = AddTagToItemByIsbn(id, newTag)
			} else {
				fmt.Println("Invalid ID type")
			}
			if err != nil {
				fmt.Printf("Error adding tag to ID %s: %v\n", id, err)
			}
			time.Sleep(500 * time.Millisecond)
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading ID file:", err)
		}
	},
}

func init() {
	updateLocationCmd.Flags().StringVarP(&location, "location", "l", "", "Location tag (e.g., Tokyo)")
	updateLocationCmd.Flags().StringVarP(&readFile, "file", "f", "", "Path to the file containing ISBNs/Item IDs")
	updateLocationCmd.Flags().StringVarP(&idType, "id", "i", "itemid", "Type of ID provided in the file (isbn or itemid)")
	updateLocationCmd.MarkFlagRequired("location")
	updateLocationCmd.MarkFlagRequired("file")
	rootCmd.AddCommand(updateLocationCmd)

	addTagCmd.Flags().StringVarP(&newTag, "tag", "t", "", "New tag to add")
	addTagCmd.Flags().StringVarP(&readFile, "file", "f", "", "Path to the file containing ISBNs/Item IDs")
	addTagCmd.Flags().StringVarP(&idType, "id", "i", "itemid", "Type of ID provided in the file (isbn or itemid)")
	addTagCmd.MarkFlagRequired("tag")
	addTagCmd.MarkFlagRequired("file")
	rootCmd.AddCommand(addTagCmd)
}

func main() {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting user home directory: %v\n", err)
		os.Exit(1)
	}

	configFile := filepath.Join(homeDir, ".config", "booklog-tool", "config.json")
	config, err = LoadConfig(configFile)
	if err != nil {
		fmt.Printf("Error loading config file: %v\n", err)
		os.Exit(1)
	}

	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
