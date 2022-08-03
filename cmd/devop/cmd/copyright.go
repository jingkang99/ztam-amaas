package cmd

import (
    "os"
    "fmt"
	"sort"
	"bytes"
	"bufio"
    "regexp"
	"strings"
	"io/ioutil"
    "path/filepath"

	"github.com/spf13/cobra"
)

var copyrightCmd = &cobra.Command{
	Use:   "copyright",
	Short: "Remove or add copyright in soure code",
	Long:  `Go through all source code files to handle copyright claim

devop.exe copyright -f list -p
devop.exe copyright -f list -a copyright.txt	
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fext, err := cmd.Flags().GetString("extension")

		if( err != nil ) {
            cmd.Help()
            os.Exit(2)
        }
        pdir, _ := cmd.Flags().GetString("dir")

		fnam, _ := cmd.Flags().GetString("file")
		regx, _ := regexp.Compile(fnam)

		fcpr, _ := cmd.Flags().GetString("addc")
		cdel, _ := cmd.Flags().GetBool("purge")

		if cdel && len(fcpr) > 0 { // either add or purge
			cmd.Help()
			os.Exit(5)
		}

		var icpr []byte
		if len(fcpr) > 0 {
			icpr, err = ioutil.ReadFile(fcpr)
			if err != nil {
				cmd.Help()
				os.Exit(5)
			}
		}

		var fcoun, fsize int64
		m := make(map[string]int64)  // file/size map

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
						m[path] = info.Size()
					}
                }else{	// search by extension
					if filepath.Ext(path) == "." + fext {
						m[path] = info.Size()

						fcoun = fcoun + 1
						fsize = fsize + info.Size()
					}
				}
                return err
            })

        if err != nil {
            fmt.Println(err)
        }

		keys := make([]string, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
		
			if ! cdel && len(fcpr) == 0 { // no add or purge
				fmt.Println(k)
				continue
			}

			readFile, err := os.Open( k )
			if err != nil {
				fmt.Println(err)
			}
			defer readFile.Close()
			
			var bs []byte
			buf := bytes.NewBuffer(bs)
			
			if len(fcpr) > 0 {	// add © info defined in txt
				buf.Write( icpr )
			}

			fScanner := bufio.NewScanner(readFile)
			fScanner.Split(bufio.ScanLines)

			imp_s, imp_e, l_num, s_blk, e_blk, bcomm := 0, 0, 0, 1, 1, false
			for fScanner.Scan() {
				line := fScanner.Text()
				lstd := strings.TrimSpace(line)

				l_num = l_num + 1

				if imp_s == 0 && imp_e == 0 { // search comment block start
					matched, _ := regexp.MatchString(`^/\*.*`, lstd)  // style /*
					if matched {
						imp_s = l_num
					}else{
						matched, _ := regexp.MatchString(`\S+`, line)
						if  s_blk  == l_num && ! matched {
							s_blk = s_blk + 1
							continue
						}

						if found, _ := regexp.MatchString(`^//`, lstd); found { // style //
							s_blk = s_blk + 1
							bcomm =  true
							continue
						}

						buf.Write( []byte(line + "\n") )
					}
				}else if imp_s > 0 && imp_e == 0 {
					matched, _ := regexp.MatchString(`.*\*/\s*`, lstd)  // style */
					if matched {  // comment block end
						imp_e = l_num
					}
				}else{	// imp_s > 0 && imp_e > 0
					matched, _ := regexp.MatchString(`\S+`, line)
					if  imp_e + e_blk == l_num && ! matched {
						e_blk = e_blk + 1
						continue
					}
					buf.Write( []byte(line + "\n") )
				}
			}
			if err := fScanner.Err(); err != nil {
				cmd.Help()
				os.Exit(5)
			}

			if imp_s > 0 || e_blk > 1 || bcomm {  // check change or not
				if  cdel {	// flag to overwrite src
					os.WriteFile(k, buf.Bytes(), 0666)
					fmt.Printf("© start %d end %d blank %d %s\n", imp_s, imp_e, e_blk, k) 
				}
			}
			
			if len(fcpr) > 0 {	// add © info defined in txt
				os.WriteFile(k, buf.Bytes(), 0666)
				fmt.Printf("© added %s\n", k) 
			}
		}

		if fnam == "" {
			fmt.Printf("\nfile count: %d\n", fcoun)
			fmt.Printf(  "total byte: %d\n", fsize)
		}
    },
}

func init() {
    copyrightCmd.PersistentFlags().StringP("extension", "e", "go", "possible values: go, js, yml")
    copyrightCmd.PersistentFlags().StringP("dir", "r", ".", "specify a project directory to scan")
	copyrightCmd.PersistentFlags().StringP("file", "f", "", "search files with a regex pattern, *.go not a proper regex")
	
	copyrightCmd.PersistentFlags().StringP("addc", "a", "", "add copyright to source files")
	copyrightCmd.PersistentFlags().BoolP("purge", "p", false, "purge copyright in source files")

	rootCmd.AddCommand(copyrightCmd)
}
