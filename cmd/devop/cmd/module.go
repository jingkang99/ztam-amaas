package cmd

import (
    "os"
    "fmt"
	"sort"
	"bufio"
	"regexp"
    "strings"
    "path/filepath"

    "github.com/spf13/cobra"
)

var moduleCmd = &cobra.Command{
    Use:   "module",
    Short: "Search directly imported Go modules",
    Long:  `Go through all .go files to parse imported modules`,
    Run: func(cmd *cobra.Command, args []string) {
        if(Debug) {
            fmt.Println("module called")
            fmt.Println("cmd args: " + strings.Join(args, ","))
        }

        mtype, err := cmd.Flags().GetString("type")
		ofile, err := cmd.Flags().GetString("save")
		exclu, err := cmd.Flags().GetString("exclude")
		count, _   := cmd.Flags().GetBool("count")
		
        if( err != nil || ( mtype != "3rd" && mtype != "sys")) {
            cmd.Help()
            os.Exit(2)
        }
        mdir, _ := cmd.Flags().GetString("dir")

		m := make(map[string]int)  // 3rd party module map
		g := make(map[string]int)  // builtin module map
		var fcoun int64

        err = filepath.Walk(mdir,
            func(path string, info os.FileInfo, err error) error {
                if err != nil {
                    return err
                }
                if filepath.Ext(path) == ".go" {
					fcoun = fcoun + 1

					readFile, err := os.Open( path )
					if err != nil {
						fmt.Println(err)
					}
					fileScanner := bufio.NewScanner(readFile)
					fileScanner.Split(bufio.ScanLines)

					imp_s, imp_e, l_num := 0, 0, 0
					for fileScanner.Scan() {
						line := fileScanner.Text()
						l_num = l_num + 1

						if imp_s == 0 { // find import start
							matched, _ := regexp.MatchString(`^import\s+\(`, line)
							if matched {
								imp_s = l_num
							}else{
								continue
							}
						}else{
							matched, _ := regexp.MatchString(`^\)$`, line)
							if matched {  // import end
								imp_e = l_num
								break
							}else{
								found, _ := regexp.MatchString(exclu, line)
								if len(exclu) > 0 && found { // exclude given pattern
									continue
								}

								matched, _ := regexp.MatchString(`"\S+\.\S+/`, line)
								if matched {	// github modules
									mdl := strings.TrimSpace(line)
									mdl  = strings.ReplaceAll(mdl, `"`, "")
									
									if strings.Index(mdl, " ") > 0 {
										s := strings.Split(mdl, " ")
										mdl = s[1]
									}
									
									matched, _ := regexp.MatchString(`^[0-9A-Za-z]+`, mdl)
									if matched { 
										m[ mdl ] = m[ mdl ] + 1
									}

									if Debug { fmt.Println("  " + mdl) }
								}else{		// buildtin modules
									mdl := strings.TrimSpace(line)
									mdl  = strings.ReplaceAll(mdl, `"`, "")
									
									matched, _ := regexp.MatchString(`^[0-9A-Za-z]+`, mdl)
									if matched { 
										g[ mdl ] = g[ mdl ] + 1
									}
								}
							}
						}
					}
					if Debug { fmt.Printf("line# %-5d %-5d %s\n", imp_s, imp_e, path) }
					
					readFile.Close()
                }

                return err
            })
			
			if mtype == "3rd" { // print sorted modules
				keys := make([]string, 0, len(m))
				for k := range m {
					keys = append(keys, k)
				}
				sort.Strings(keys)

				f, _ := os.Create(ofile)
				defer f.Close()
				w := bufio.NewWriter(f)

				for _, k := range keys {			
					if count {
						fmt.Printf("%-10d %s\n", m[k], k)
					}else{
						fmt.Println(k)
						w.WriteString(k + "\n")
					}
				}
				w.Flush()
			}else{
				keys := make([]string, 0, len(g))
				for k := range g {
					keys = append(keys, k)
				}
				sort.Strings(keys)

				for _, k := range keys {
					if count {
						fmt.Printf("%-10d %s\n", g[k], k)
					}else{
						fmt.Println(k)
					}
				}
			}

			if count {
				fmt.Printf("\nfile count : %d\n", fcoun)
				fmt.Printf("sys  module: %d\n", len(g))
				fmt.Printf("3rd  module: %d\n", len(m))
			}

        if err != nil {
            fmt.Println(err)
        }
    },
}

func init() {
    moduleCmd.PersistentFlags().StringP("type", "t", "3rd", "possible values: 3rd, sys")
    moduleCmd.PersistentFlags().StringP("dir",  "r", ".", "specify a project directory to scan")
	moduleCmd.PersistentFlags().StringP("save", "o", "module.txt", "save result to a file")
	moduleCmd.PersistentFlags().StringP("exclude", "x", "", "exlcude files with a regex pattern")

    moduleCmd.PersistentFlags().BoolP("count", "n", false, "print file and module count")
	

    listCmd.AddCommand(moduleCmd)
}

func CheckFile() {

}

