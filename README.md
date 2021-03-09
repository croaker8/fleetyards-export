# **NOTE** This program will currently not work as third party logins are currently disabled on fleetyards.net

# fleetyards-export
Tool for exporting your public hanger from https://fleetyards.net.

This command will access your public hanger URL of fleetyards.net and export the data for your vehicles to a CSV file (default is output.csv).

Edit the export-field-list file and add # at beginning of line to exclude specific fields from the CSV output.

## Build
    go build -o fleetyards-export
## Usage
Usage is shown if -h (or --help) flag is supplied. The -u (or --user) flag is required. Password may be specificed using the -p (or --pass) flag or if not specified on the command line the user will be prompted to enter the password and the input will be hidden.

    fleetyards-export is a tool for exporting public hanger data from fleetyards.net

    Usage:
      fleetyards-export [flags]

    Flags:
      -f, --fields string   Path to fields list file (default "export-field-list")
      -h, --help            help for fleetyards-export
      -o, --output string   Path to output file (default "output.csv")
      -u, --user string     User name on fleetyards.net (required)
