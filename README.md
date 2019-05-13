# fleetyards-export
Tool for exporting hanger data from https://fleetyards.net.

This command will login to your fleetyards.net account and export the data for the vehicles in your hanger to a CSV file (default is output.csv).

Edit the export-field-list file and add # at beginning of line to exclude specific fields from the CSV output.

## Usage
Usage show if -h or --help flag is supplied. The -u (or --user) flag is required. Password may be specificed using the -p (or --pass) flag or if not specified on the command line the user will be prompted to enter the password and the input will be hidden.

    fleetyards-export is a tool for exporting hanger data from fleetyards.net

    Usage:
      fleetyards-export [flags]

    Flags:
      -f, --fields string   Path to fields list file (default "export-field-list")
      -h, --help            help for fleetyards-export
      -o, --output string   Path to output file (default "output.csv")
      -p, --pass string     Password to login to fleetyards.net
      -u, --user string     User name to login to fleetyards.net (required)
