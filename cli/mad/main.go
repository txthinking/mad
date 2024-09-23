package main

import (
	"errors"
	"log"
	"net"
	"os"
	"time"

	"github.com/txthinking/mad"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "mad"
	app.Version = "20240923"
	app.Usage = "Generate root CA and derivative certificate for any domains and any IPs"
	app.Authors = []*cli.Author{
		{
			Name:  "Cloud",
			Email: "cloud@txthinking.com",
		},
	}
	app.Commands = []*cli.Command{
		&cli.Command{
			Name:  "ca",
			Usage: "Generate CA",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "ca",
					Usage: "CA file which will be created or overwritten",
					Value: "ca.pem",
				},
				&cli.StringFlag{
					Name:  "key",
					Usage: "Key file which will be created or overwritten",
					Value: "ca.key.pem",
				},
				&cli.StringFlag{
					Name:  "organization",
					Value: "github.com/txthinking/mad",
				},
				&cli.StringFlag{
					Name:  "organizationUnit",
					Value: "github.com/txthinking/mad",
				},
				&cli.StringFlag{
					Name:  "commonName",
					Value: "github.com/txthinking/mad",
				},
				&cli.StringFlag{
					Name:  "start",
					Usage: "Certificate valid start time, such as: '2024-09-22T13:07:38+08:00'. If empty, it is the current time",
				},
				&cli.StringFlag{
					Name:  "end",
					Usage: "Certificate valid end time, such as: '2024-09-22T13:07:38+08:00'. If empty, it is start time add 10 years",
				},
				&cli.BoolFlag{
					Name:  "install",
					Usage: "Install immediately after creation",
				},
			},
			Action: func(c *cli.Context) error {
				var err error
				start := time.Now()
				if c.String("start") != "" {
					start, err = time.Parse(time.RFC3339, c.String("start"))
					if err != nil {
						return err
					}
				}
				end := start.AddDate(10, 0, 0)
				if c.String("end") != "" {
					end, err = time.Parse(time.RFC3339, c.String("end"))
					if err != nil {
						return err
					}
				}
				ca := mad.NewCa(c.String("organization"), c.String("organizationUnit"), c.String("commonName"), start, end)
				if err := ca.Create(); err != nil {
					return err
				}
				if err := ca.SaveToFile(c.String("ca"), c.String("key")); err != nil {
					return err
				}
				if c.Bool("install") {
					if err := mad.Install(c.String("ca")); err != nil {
						return err
					}
				}
				return nil
			},
		},
		&cli.Command{
			Name:  "cert",
			Usage: "Generate certificate",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "ca",
					Usage: "ROOT CA file path",
					Value: "ca.pem",
				},
				&cli.StringFlag{
					Name:  "ca_key",
					Usage: "Deprecated, please use --caKey",
				},
				&cli.StringFlag{
					Name:  "caKey",
					Usage: "ROOT Key file path",
					Value: "ca.key.pem",
				},
				&cli.StringFlag{
					Name:  "cert",
					Usage: "Certificate file which will be created or overwritten",
					Value: "cert.pem",
				},
				&cli.StringFlag{
					Name:  "key",
					Usage: "Certificate key file which will be created or overwritten",
					Value: "cert.key.pem",
				},
				&cli.StringFlag{
					Name:  "organization",
					Value: "github.com/txthinking/mad",
				},
				&cli.StringFlag{
					Name:  "organizationUnit",
					Value: "github.com/txthinking/mad",
				},
				&cli.StringSliceFlag{
					Name:  "ip",
					Usage: "IP address",
				},
				&cli.StringSliceFlag{
					Name:  "domain",
					Usage: "Domain name",
				},
				&cli.StringFlag{
					Name:  "commonName",
					Usage: "If empty, the first domain or IP will be used",
				},
				&cli.StringFlag{
					Name:  "start",
					Usage: "Certificate valid start time, such as: '2024-09-22T13:07:38+08:00'. If empty, it is the current time",
				},
				&cli.StringFlag{
					Name:  "end",
					Usage: "Certificate valid end time, such as: '2024-09-22T13:07:38+08:00'. If empty, it is start time add 10 years",
				},
			},
			Action: func(c *cli.Context) error {
				ca, err := os.ReadFile(c.String("ca"))
				if err != nil {
					return err
				}
				var caKey []byte
				if c.String("ca_key") != "" {
					caKey, err = os.ReadFile(c.String("ca_key"))
					if err != nil {
						return err
					}
				} else {
					caKey, err = os.ReadFile(c.String("caKey"))
					if err != nil {
						return err
					}
				}
				start := time.Now()
				if c.String("start") != "" {
					start, err = time.Parse(time.RFC3339, c.String("start"))
					if err != nil {
						return err
					}
				}
				end := start.AddDate(10, 0, 0)
				if c.String("end") != "" {
					end, err = time.Parse(time.RFC3339, c.String("end"))
					if err != nil {
						return err
					}
				}
				cert := mad.NewCert(ca, caKey, c.String("organization"), c.String("organizationUnit"), start, end)
				ips := make([]net.IP, 0)
				for _, v := range c.StringSlice("ip") {
					ip := net.ParseIP(v)
					if ip == nil {
						return errors.New(v + " is not an IP")
					}
					ips = append(ips, ip)
				}
				if len(ips) > 0 {
					cert.SetIPAddresses(ips)
				}
				if len(c.StringSlice("domain")) > 0 {
					cert.SetDNSNames(c.StringSlice("domain"))
				}
				if c.String("commonName") != "" {
					cert.SetCommonName(c.String("commonName"))
				}
				if err := cert.Create(); err != nil {
					return err
				}
				if err := cert.SaveToFile(c.String("cert"), c.String("key")); err != nil {
					return err
				}
				return nil
			},
		},
		&cli.Command{
			Name:  "install",
			Usage: "Install ROOT CA",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "ca",
					Usage: "CA file which will be installed",
					Value: "ca.pem",
				},
			},
			Action: func(c *cli.Context) error {
				if err := mad.Install(c.String("ca")); err != nil {
					return err
				}
				return nil
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
