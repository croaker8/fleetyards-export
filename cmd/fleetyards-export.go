package cmd

import (
	"fmt"
	"github.com/quarkstar8/fleetyards-export/fleetyards"
	"github.com/spf13/cobra"
	"os"
)

var flyExCmd = &cobra.Command{
	Use:   "fleetyards-export",
	Short: "fleetyards-export is a tool for exporting public hanger data from fleetyards.net",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		// promt for password if arg not set
		//if password == "" {
		//	fmt.Print("Enter Password: ")
		//	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		//	if err != nil {
		//		fmt.Printf("Error reading password from terminal: %s", err)
		//		os.Exit(1)
		//	} else {
		//		fmt.Println()
		//	}
		//	password = string(bytePassword)
		//}

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
		//fmt.Println("Connecting to service...")
		//token, err := fleetyards.Signin(username, password)
		//if err != nil {
		//	os.Exit(6)
		//}

		// get the data and export
		fmt.Println("Exporting data...")
		exportErr := fleetyards.ExportHangerToCsv(username, outputPath, fieldList)

		// signout from service
		//fmt.Println("Disconnecting from service...")
		//err = fleetyards.Signout(token)
		//if err != nil {
		//	os.Exit(7)
		//}

		// set error code signout worked but export failed
		if exportErr != nil {
			os.Exit(8)
		}

		fmt.Println("Export complete.")
	},
}

var username string
var fieldsPath string
var outputPath string

func init() {
	flyExCmd.Flags().StringVarP(&username, "user", "u", "", "User name on fleetyards.net (required)")
	flyExCmd.Flags().StringVarP(&fieldsPath, "fields", "f", "export-field-list", "Path to fields list file")
	flyExCmd.Flags().StringVarP(&outputPath, "output", "o", "output.csv", "Path to output file")

	flyExCmd.MarkFlagRequired("user")
}

// Execute -- execute the command
func Execute() {
	if err := flyExCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
