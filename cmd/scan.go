/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/Adit0507/pScan.com/scan"
	"github.com/spf13/cobra"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Run a port scan on the hosts",
	RunE: func(cmd *cobra.Command, args []string) error {
		hostsFile, err := cmd.Flags().GetString("hosts-file")
		if err != nil {
			return err
		}
		
		ports, err := cmd.Flags().GetIntSlice("ports")
		if err != nil {
			return err
		}

		return scanAction(os.Stdout, hostsFile, ports)
	},
}

func scanAction(out io.Writer, hostsFile string, ports []int) error {
	hl := &scan.HostLists{}

	if err := hl.Load(hostsFile); err != nil {
		return err
	}

	res := scan.Run(hl, ports)

	return printResults(out, res)
}

func printResults(out io.Writer, res []scan.Results) error {
	msg := ""

	for _, r := range res {
		msg += fmt.Sprintf("%s:", r.Host)
	
		if r.NotFound {
			msg += fmt.Sprintf("Host not found\n\n")
			continue
		}

		msg += fmt.Sprintln()

		for _, p := range r.PortStates {
			msg += fmt.Sprintf("\t %d: %s \n " , p.Port, p.Open)
		}
		
		msg += fmt.Sprintln()
	}

	_, err := fmt.Fprintf(out, msg)
	return err
}

func init() {
	rootCmd.AddCommand(scanCmd)

	scanCmd.Flags().IntSliceP("ports", "p", []int{22,80, 443}, "ports to scan")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
