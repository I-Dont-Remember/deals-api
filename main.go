package main

import (
    "os"
    "github.com/urfave/cli"
    helper "github.com/I-Dont-Remember/deals-api/awshelper"
    tools "github.com/I-Dont-Remember/deals-api/tools"
)

var (
    // assume local is true for now
    local = true
)

// getCommands returns the full list of commands for the CLI tool
func getCommands() []cli.Command {
    return []cli.Command {
        {
            Name: "tables",
            Usage: "options for tables",
            Subcommands: []cli.Command {
                {
                    Name: "create",
                    Usage: "create table from predefined list",
                    Flags: []cli.Flag {
                        cli.BoolFlag{
                            Name: "all, a",
                            Usage: "Create all predefined tables",
                        },
                    },

                    Action: func(c *cli.Context) error {
                        if c.Bool("all") {
                            for key,_ := range helper.Tables {
                                err := helper.CreateTable(local, key)
                                tools.Check(err)
                            }
                        } else {
                            name := c.Args().First()
                            err := helper.CreateTable(local, name)
                            tools.Check(err)
                        }
                        return nil
                    },
                },
            },
        },
    }
}

func main() {
    app := cli.NewApp()
    app.Name = "cli"
    app.Usage = "cli for Madtown Deals repo"

    app.Commands = getCommands()

    app.Run(os.Args)
}
