package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var (
	protoFilePath string // 源 .proto 文件的路径
)

var pbCmd = &cobra.Command{
	Use:   "pb",
	Short: "编译 proto 文件生成 Go 代码",
	Run: func(cmd *cobra.Command, args []string) {
		err := exec.Command(
			"protoc",
			"--go_out=./",
			"--go_opt=paths=source_relative",
			"--go-grpc_out=./",
			"--go-grpc_opt=paths=source_relative",
			protoFilePath,
		).Run()
		if err != nil {
			return
		}
		fmt.Println("proto 编译成功！")
	},
}

func init() {
	rootCmd.AddCommand(pbCmd)
	pbCmd.Flags().StringVarP(&protoFilePath, "file", "f", "", "源 .proto 文件的路径（必填）")
	_ = pbCmd.MarkFlagRequired("file")
}
