package main

import (
    "fmt"
    "os"
    "github.com/urfave/cli"
    helper "github.com/I-Dont-Remember/gimmedeals/api/awshelper"
)

var (
    // assume local is true for now
    local = true
)

// check consolidates a common error checking mechanism into one call
func check(err error) {
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }
}

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
                                check(err)
                            }
                        } else {
                            name := c.Args().First()
                            err := helper.CreateTable(local, name)
                           check(err)
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
    app.Name = "tools"
    app.Usage = "Tools for Madtown Deals repo"

    app.Commands = getCommands()

    app.Run(os.Args)
}
