/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bookapi/routes"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long:  `A longer description that spans multiple lines and likely contains examples`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚀 Starting server on http://127.0.0.1:8080...")

		// Use your actual Routes() function to get the router
		r := routes.Routes()

		if err := http.ListenAndServe(":8080", r); err != nil {
			fmt.Printf("❌ Server failed: %s\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
