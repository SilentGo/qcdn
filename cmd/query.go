package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"qcdn/modules"
)

var (
	host string
)

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "A brief description of your command",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		c := modules.NewClient(secretID, secretKey, host, 0, 0, 500)
		query(c, host)
	},
}

func init() {
	rootCmd.AddCommand(queryCmd)
	queryCmd.Flags().StringVarP(&host, "domain", "d", "", "")
}

func query(c *modules.Client, host string) {
	if host != "" {
		rsp := c.GetHostInfoByHost()
		printRsp(rsp)
	} else {
		rsp := c.DescribeCdnHosts()
		printRsp(rsp)
	}
}

func printRsp(rsp *modules.DescribeCdnHostsRsp) {
	if rsp.Code == 0 {
		for i := 0; i < len(rsp.Data.Hosts); i++ {
			fmt.Printf("域名：%s  记录： %s\n", rsp.Data.Hosts[i].Host, rsp.Data.Hosts[i].Cname)
		}
	} else {
		fmt.Println(rsp.Message)
	}
}
