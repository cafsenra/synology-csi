/*
 * Copyright 2021 Synology Inc.
 */
package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"github.com/SynologyOpenSource/synology-csi/pkg/dsm/common"
	"github.com/SynologyOpenSource/synology-csi/pkg/dsm/webapi"
)

var https = false
var port = -1

var cmdDsm = &cobra.Command{
	Use:   "dsm",
	Short: "dsm",
	Long:  `dsm`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var cmdDsmLogin = &cobra.Command{
	Use:   "login <ip> <username> <password>",
	Short: "login dsm",
	Args:  cobra.MinimumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		var defaultPort = 5000
		if https {
			defaultPort = 5001
		}
		if port != -1 {
			defaultPort = port
		}

		dsmApi := &webapi.DSM{
			Ip:       args[0],
			Username: args[1],
			Password: args[2],
			Port:     defaultPort,
			Https:    https,
		}

		err := dsmApi.Login()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

// Always get the first client from ClientInfo for synocli testing
func LoginDsmForTest() (*webapi.DSM, error) {
	info, err := common.LoadConfig("./config/client-info.yml")
	if err != nil {
		return nil, fmt.Errorf("Failed to read config: %v", err)
	}
	if len(info.Clients) == 0 {
		return nil, fmt.Errorf("No client in client-info.yml")
	}

	dsm := &webapi.DSM{
		Ip:       info.Clients[0].Host,
		Port:     info.Clients[0].Port,
		Username: info.Clients[0].Username,
		Password: info.Clients[0].Password,
		Https:    info.Clients[0].Https,
	}

	if err := dsm.Login(); err != nil {
		return nil, fmt.Errorf("Failed to login to DSM: [%s]. err: %v", dsm.Ip, err)
	}
	return dsm, nil
}

func init() {
	cmdDsm.AddCommand(cmdDsmLogin)

	cmdDsmLogin.PersistentFlags().BoolVar(&https, "https", false, "Use HTTPS to login DSM")
	cmdDsmLogin.PersistentFlags().IntVarP(&port, "port", "p", -1, "Use assigned port to login DSM")
}