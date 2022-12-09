package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	
	"github.com/jingkang99/ztam-amaas/pkg/global"
)

var Debug bool

var rootCmd = &cobra.Command{
	Use:   "devop",
	Short: "Scan project basic info",
	Long:  ``,

}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&global.CFile, "config", "c", "", "specify config file")
	rootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", false, "debug mode, print more info")

	global.Debug = Debug
	
}

func initConfig() {

	if global.CFile != "" {
		viper.SetConfigFile(global.CFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.AddConfigPath(".")

		viper.SetConfigName("devop")
		viper.SetConfigType("toml")	// ini not effective if toml exists
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if(Debug){
			fmt.Println("Error: " + viper.ConfigFileUsed(), err)
		}
	}

	var conf = &global.CFObj
	
	if err := viper.Unmarshal(conf); err != nil {
		if(Debug){
			fmt.Printf("Error: cannot read: %s", err)
		}
	}

	if(Debug){
		fmt.Printf("%+v\n", global.CFObj)
		fmt.Fprintln(os.Stderr, "\nconf:", viper.ConfigFileUsed())
		viper.Debug()
	}
}