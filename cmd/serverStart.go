/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/AkronimBlack/process-manager/pkg/parser"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"log"
)

var (
	router             *gin.Engine
	serverFileLocation string
)

// serverStartCmd represents the serverStart command
var serverStartCmd = &cobra.Command{
	Use:   "server:start",
	Short: "Start a process server",
	Run: func(cmd *cobra.Command, args []string) {
		spinUp()
	},
}

func init() {
	rootCmd.AddCommand(serverStartCmd)
	serverStartCmd.Flags().StringVarP(&serverFileLocation, "file-location", "f", "", "location of json file to parse")
}

func spinUp() {
	router = gin.Default()
	parser.BuildHttp(router, serverFileLocation)
	err := router.Run("0.0.0.0:8080")
	if err != nil {
		log.Panic(err)
	}
}
