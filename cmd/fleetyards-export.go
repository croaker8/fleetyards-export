package cmd

import (
	"fmt"
	"os"

	"github.com/croaker8/fleetyards-export/fleetyards"
	"github.com/spf13/cobra"
)

var flyExCmd = &cobra.Command{
	Use:   "fleetyards-export",
	Short: "fleetyards-export is a tool for exporting hanger data from fleetyards.net",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		// get the field list to export
		fieldList, err := getFieldList(fieldsPath)
		if err != nil {
			os.Exit(2)
		}

		// remove any existing output file
		stat, err := os.Stat(outputPath)
		if err != nil && !os.IsNotExist(err) {
			fmt.Printf("Error getting stat of existing output file: %s\n", err)
			os.Exit(3)
		}

		// check output is a directory
		if stat != nil && stat.IsDir() {
			fmt.Printf("Error, output path '%s' is a directory\n", outputPath)
			os.Exit(4)
		}

		// if output file exists remove it
		if stat != nil {
			err = os.Remove(outputPath)
			if err != nil {
				fmt.Printf("Error deleting existing output file '%s': %s\n", outputPath, err)
				os.Exit(5)
			}
		}

		// signin to fleetyards service
		token, err := fleetyards.Signin(username, password)
		if err != nil {
			os.Exit(6)
		}

		// get the data and export
		exportErr := fleetyards.ExportHangerToCsv(token, outputPath, fieldList)

		// signout from service
		err = fleetyards.Signout(token)
		if err != nil {
			os.Exit(7)
		}

		// set error code signout worked but export failed
		if exportErr != nil {
			os.Exit(8)
		}

	},
}

var username string
var password string
var fieldsPath string
var outputPath string

func init() {
	flyExCmd.Flags().StringVarP(&username, "user", "u", "", "User name to login to fleetyards.net")
	flyExCmd.Flags().StringVarP(&password, "pass", "p", "", "Password to login to fleetyards.net")
	flyExCmd.Flags().StringVarP(&fieldsPath, "fields", "f", "export-field-list", "Path to fields list file")
	flyExCmd.Flags().StringVarP(&outputPath, "output", "o", "output.csv", "Path to output file")

	flyExCmd.MarkFlagRequired("user")
	flyExCmd.MarkFlagRequired("pass")
}

// Execute -- execute the command
func Execute() {
	if err := flyExCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
