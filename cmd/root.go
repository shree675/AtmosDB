package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"atmosdb/server"
	"atmosdb/util"
)

var rootCmd = &cobra.Command{
	Use:   "atmos",
	Short: "A simple concurrent in-memory database",
	Run: func(cmd *cobra.Command, args []string) {
		printAscii()
		server.Start()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func printAscii() {
	util.PrintBlue(`
 _______                        ______  ______  
(_______)  _                   (______)(____  \ 
 _______ _| |_ ____   ___   ___ _     _ ____)  )
|  ___  (_   _)    \ / _ \ /___) |   | |  __  ( 
| |   | | | |_| | | | |_| |___ | |__/ /| |__)  )
|_|   |_|  \__)_|_|_|\___/(___/|_____/ |______/ `)
	fmt.Println()
}
