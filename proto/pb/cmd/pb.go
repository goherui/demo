/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os/exec"

	"github.com/spf13/cobra"
)

var (
	protoFile string
	outputDir string
)

// pbCmd represents the pb command
var pbCmd = &cobra.Command{
	Use:   "pb",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := exec.Command(
			"protoc",
			"--go_out=./"+outputDir,
			"--go_opt=paths=source_relative",
			"--go-grpc_out=./"+outputDir,
			"--go-grpc_opt=paths=source_relative",
			"cmd/goods.proto",
		).Run()
		if err != nil {
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(pbCmd)
	pbCmd.Flags().StringVarP(&protoFile, "file", "f", "", "proto文件路径")
	pbCmd.Flags().StringVarP(&protoFile, "dir", "d", "", "pb文件生成路径")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pbCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pbCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
