package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	folder string

	readfile string

	fileExt string

	rootCmd = &cobra.Command{
		Use:   "generateQRCode",
		Short: "Generate QR Code for your pin Code.",
		RunE: func(cmd *cobra.Command, args []string) error {
			r, err := checkFlags()
			if err != nil {
				errMsg(err)
				return err
			}
			return r.generateQRCode()
		},
	}
)

type flags struct {
	folder string

	readfile string

	fileExt string
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&readfile, "readfile", "r", "pinCodeFile.txt", "read file content for pin code (default pinCodeFile.txt)")
	rootCmd.PersistentFlags().StringVarP(&folder, "folder", "f", "file", "create folder (default file)")
	rootCmd.PersistentFlags().StringVarP(&fileExt, "fileExt", "e", ".png", "file extension (default .png)")
}

func checkFlags() (fs *flags, err error) {
	if len(readfile) == 0 {
		return nil, fmt.Errorf("flags readfile is empty")
	}

	if len(folder) == 0 {
		return nil, fmt.Errorf("flags folder is empty")
	}

	if len(fileExt) == 0 {
		return nil, fmt.Errorf("flags fileExt is empty")
	}

	fs = &flags{
		readfile: readfile,
		folder:   folder,
		fileExt:  fileExt,
	}

	return fs, nil
}

func errMsg(msg interface{}) {
	fmt.Println("Error:", msg)
}

// Execute command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
