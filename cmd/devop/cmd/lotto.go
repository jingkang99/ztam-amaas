package cmd

import (
	"os"
	"fmt"
	"math/rand"

	"github.com/spf13/cobra"

	"github.com/jingkang99/ztam-amaas/pkg/global"
)

// lottoCmd represents the lotto command
var lottoCmd = &cobra.Command{
	Use:   "lotto",
	Short: "Pick lottery numbers",
	Long:  `Pick lottery numbers randomly

Support: "mega-millions", "powerball", "superlotto-plus", "fantasy-5", "lotto-test"

devop lotto mega-millions -n 5
devop lotto lotto-test -m 47 -p 5 -a 27 -t 1 -n 10`,
	Args:  cobra.OnlyValidArgs,
	ValidArgs: []string{"mega-millions", "powerball", "superlotto-plus", "fantasy-5", "lotto-test"},
	PreRunE: func(cmd *cobra.Command, args []string) error {
        if len(args) != 1 {
            cmd.Help()
            os.Exit(0)
        }
        return nil
    },
	Run: func(cmd *cobra.Command, args []string) {
		for i := 0; i < COUNT; i++ {
			switch args[0] {
			case "mega-millions":
				lotto :=  global.Lotto(70, 5)

				var cSeed global.CryptoSrc
				mega := rand.New(cSeed)

				fmt.Printf("%2v %s%d%s\n", lotto, string(global.ColorBlue), mega.Intn(25) + 1, string(global.ColorReset))
				
			case "powerball":
				lotto :=  global.Lotto(69, 5)

				var cSeed global.CryptoSrc
				mega := rand.New(cSeed)

				//fmt.Printf("%2v %d\n", lotto, mega.Intn(26) + 1)
				fmt.Printf("%2v %s%d%s\n", lotto, string(global.ColorGreen), mega.Intn(26) + 1, string(global.ColorReset))

			case "superlotto-plus":
				lotto :=  global.Lotto(47, 5)

				var cSeed global.CryptoSrc
				mega := rand.New(cSeed)

				fmt.Printf("%2v %s%d%s\n", lotto, string(global.ColorCyan), mega.Intn(27) + 1, string(global.ColorReset))

			case "fantasy-5":
				lotto :=  global.Lotto(39, 5)
				fmt.Printf("%s%2v%s\n", string(global.ColorRed), lotto, string(global.ColorReset))

			case "lotto-test":
				
				lotto :=  global.Lotto(MAXA, NUMA)
				megan :=  global.Lotto(MAXM, NUMM)

				fmt.Printf("%s%2v%s  ", string(global.ColorRed),    lotto, string(global.ColorReset))
				fmt.Printf("%s%2v%s\n", string(global.ColorPurple), megan, string(global.ColorReset))
			}
		}
	},
}

var COUNT, MAXA, NUMA, MAXM, NUMM int

func init() {
	lottoCmd.PersistentFlags().IntVarP(&COUNT, "count", "n", 1, "Generate lotto count")
	
	lottoCmd.PersistentFlags().IntVarP(&MAXA, "max-numb", "m", 47, "pick range 1 to max-num")
	lottoCmd.PersistentFlags().IntVarP(&NUMA, "numb-cnt", "p", 5, "pick count")

	lottoCmd.PersistentFlags().IntVarP(&MAXM, "max-mega", "a", 27, "pick range 1 to max-mega")
	lottoCmd.PersistentFlags().IntVarP(&NUMM, "mega-cnt", "t", 1, "pick count")

	rootCmd.AddCommand(lottoCmd)
}

// awk -F']' '{print $2}' 2 | sort -n  | uniq -c
/*
1	36955	-0.221%
2	36876	-0.435%
3	37213	0.475%
4	36957	-0.216%
5	37238	0.543%
6	36848	-0.510%
7	37127	0.243%
8	37120	0.224%
9	37044	0.019%
10	36733	-0.821%
11	37221	0.497%
12	36995	-0.113%
13	37045	0.022%
14	36885	-0.410%
15	36854	-0.494%
16	37173	0.367%
17	36795	-0.653%
18	36773	-0.713%
19	37107	0.189%
20	37239	0.545%
21	37358	0.867%
22	37101	0.173%
23	37081	0.119%
24	36854	-0.494%
25	37017	-0.054%
26	37162	0.337%
27	37229	0.518%
		
	37037.03704	
*/
