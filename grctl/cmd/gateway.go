// RAINBOND, Application Management Platform
// Copyright (C) 2014-2020 Goodrain Co., Ltd.

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version. For any non-GPL usage of Rainbond,
// one or multiple Commercial Licenses authorized by Goodrain Co., Ltd.
// must be obtained first.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/goodrain/rainbond/gateway/controller/openresty/model"
	"github.com/gosuri/uitable"
	"github.com/urfave/cli"
)

//NewCmdGateway gateway cmd
func NewCmdGateway() cli.Command {
	c := cli.Command{
		Name:  "gateway",
		Usage: "Gateway management related commands",
		Subcommands: []cli.Command{
			cli.Command{
				Name:  "endpoints",
				Usage: "list gateway http endpoints",
				Flags: []cli.Flag{
					cli.IntFlag{
						Name:  "port",
						Usage: "gateway http endpoint query port",
						Value: 18080,
					},
				},
				Action: func(c *cli.Context) error {
					return listHTTPEndpoint(c)
				},
			},
		},
	}
	return c
}

func listHTTPEndpoint(c *cli.Context) error {
	res, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/config/backends", c.Int("port")))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if res.Body != nil {
		defer res.Body.Close()
		decoder := json.NewDecoder(res.Body)
		var backends []*model.Backend
		if err := decoder.Decode(&backends); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		table := uitable.New()
		table.Wrap = true // wrap columns
		for _, b := range backends {
			table.AddRow(b.Name, strings.Join(func() []string {
				var re []string
				for _, e := range b.Endpoints {
					re = append(re, fmt.Sprintf("%s:%s %d", e.Address, e.Port, e.Weight))
				}
				return re
			}(), ";"))
		}
		fmt.Println(table)
	}
	return nil
}
