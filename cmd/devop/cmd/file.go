package cmd

import (
    "os"
    "fmt"
    "regexp"
	"strings"
    "path/filepath"

	"github.com/spf13/cobra"
)

var fileCmd = &cobra.Command{
	Use:   "file",
	Short: "List files with filter condition",
	Long: `Searh files with given condition recursively`,
	Run: func(cmd *cobra.Command, args []string) {
		fext, err := cmd.Flags().GetString("extension")

		if( err != nil ) {
            cmd.Help()
            os.Exit(2)
        }
        pdir, _ := cmd.Flags().GetString("dir")
		size, _ := cmd.Flags().GetBool("size")

		fnam, _ := cmd.Flags().GetString("file")
		regx, _ := regexp.Compile(fnam)

		var fcoun, fsize int64

        if(Debug) {
			fmt.Println("-----")
            fmt.Println("file called")
            fmt.Println("file args: " + strings.Join(args, ","))
			fmt.Println("file pattern: " + fnam)
        }
		if len(fnam) > 0 && fnam[0:1] == "*"  {
            cmd.Help()
            os.Exit(4)
        }

        err = filepath.Walk(pdir,
            func(path string, info os.FileInfo, err error) error {
                if err != nil {
                    return err
                }

				if fnam != "" {	// search by file pattern
					if regx.MatchString(path) {
						if size {
							fmt.Printf("%-10d %s\n", info.Size(), path)
						}else{
							fmt.Println(path)
						}						
					}
                }else{	// search by extension
					if filepath.Ext(path) == "." + fext {
						if size {
							fmt.Printf("%-10d %s\n", info.Size(), path)
						}else{
							fmt.Println(path)
						}

						fcoun = fcoun + 1
						fsize = fsize + info.Size()
					}
				}
                return err
            })

        if err != nil {
            fmt.Println(err)
        }

		if fnam == "" {
			fmt.Printf("\nfile count: %d\n", fcoun)
			fmt.Printf(  "total byte: %d\n", fsize)
		}
    },
}

func init() {
    fileCmd.PersistentFlags().StringP("extension", "e", "go", "possible values: go, js, yml")
    fileCmd.PersistentFlags().StringP("dir", "r", ".", "specify a project directory to scan")
	fileCmd.PersistentFlags().StringP("file", "f", "", "search files with a regex pattern, *.go not a proper regex")
	fileCmd.PersistentFlags().BoolP("size", "s", false, "print file size")

	listCmd.AddCommand(fileCmd)
}
