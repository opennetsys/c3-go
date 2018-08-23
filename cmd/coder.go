package cmd

import (
	"fmt"

	"github.com/c3systems/c3-go/common/txparamcoder"
	"github.com/spf13/cobra"
)

func encodeCmd() *cobra.Command {
	var (
		data   string
		method string
	)

	encodecmd := &cobra.Command{
		Use:   "encode",
		Short: "Encode string data",
		Long:  "Encode string data to be used for transaction data",
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if method != "" {
				fmt.Println(txparamcoder.EncodeMethodName(method))
			}
			if data != "" {
				fmt.Println(txparamcoder.EncodeParam(data))
			}

			return nil
		},
	}

	encodecmd.Flags().StringVarP(&data, "data", "d", "", "Data to encode")
	encodecmd.Flags().StringVarP(&method, "method", "m", "", "Method name to encode")

	return encodecmd
}
