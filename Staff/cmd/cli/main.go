package main

import (
	"log"
	"os"
	"sort"

	"github.com/duy9911/Staff/handler"
	"github.com/duy9911/Staff/models"
	"github.com/urfave/cli/v2"
)

func main() {
	staff := models.Staff{}
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "create staff, eg:  -n <name> -d <1999-09-15> -s <153.0> create  ",
				Action: func(c *cli.Context) error {
					handler.CreateStaff(staff)
					return nil
				},
			},
			{
				Name:    "getall",
				Aliases: []string{"g"},
				Usage:   "return all staffs, eg: getall <domain> ",
				Action: func(c *cli.Context) error {
					domain := c.Args().First()
					handler.ReturnStaffs(domain)
					return nil
				},
			},
			{
				Name:    "update",
				Aliases: []string{"u"},
				Usage:   "update, eg: -flag update <staff_id> ",
				Action: func(c *cli.Context) error {
					key := c.Args().First()
					handler.UpdateStaff(key, staff)
					return nil
				},
			},
			{
				Name:    "delete",
				Aliases: []string{"d"},
				Usage:   "delete, eg: delete <staff_id>",
				Action: func(c *cli.Context) error {
					key := c.Args().First()
					handler.Deletestaff(key)
					return nil
				},
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "name",
				Aliases:     []string{"n"},
				Value:       " ",
				Usage:       "name for staff",
				Destination: &staff.Name,
			},
			&cli.StringFlag{
				Name:        "dob",
				Aliases:     []string{"d"},
				Value:       " ",
				Usage:       "day of birth for staff, following format 2006-01-02",
				Destination: &staff.Dob,
			},
			&cli.StringFlag{
				Name:        "gender",
				Aliases:     []string{"g"},
				Value:       " ",
				Usage:       "gender of birth for staff, male/female/both..",
				Destination: &staff.Gender,
			},
			&cli.Float64Flag{
				Name:        "salary",
				Aliases:     []string{"s"},
				Value:       0,
				Usage:       "salary for staff, number only (allow float)",
				Destination: &staff.Salary,
			},
		},
	}
	sort.Sort(cli.CommandsByName(app.Commands))
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
