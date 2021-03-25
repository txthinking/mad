package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net"
	_ "net/http/pprof"
	"os"

	"github.com/txthinking/mad"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "Mad"
	app.Version = "20210401"
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
					Value: "ca_key.pem",
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
				&cli.BoolFlag{
					Name:  "install",
					Usage: "Install CA",
				},
			},
			Action: func(c *cli.Context) error {
				ca := mad.NewCa(c.String("organization"), c.String("organizationUnit"), c.String("commonName"))
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
					Usage: "ROOT Key file path",
					Value: "ca_key.pem",
				},
				&cli.StringFlag{
					Name:  "cert",
					Usage: "Certificate file which will be created or overwritten",
					Value: "cert.pem",
				},
				&cli.StringFlag{
					Name:  "key",
					Usage: "Certificate key file which will be created or overwritten",
					Value: "cert_key.pem",
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
					Usage: "IP address. Repeated",
				},
				&cli.StringSliceFlag{
					Name:  "domain",
					Usage: "Domain name. Repeated",
				},
			},
			Action: func(c *cli.Context) error {
				ca, err := ioutil.ReadFile(c.String("ca"))
				if err != nil {
					return err
				}
				caKey, err := ioutil.ReadFile(c.String("ca_key"))
				if err != nil {
					return err
				}
				cert := mad.NewCert(ca, caKey, c.String("organization"), c.String("organizationUnit"))
				ips := make([]net.IP, 0)
				for _, v := range c.StringSlice("ip") {
					ip := net.ParseIP(v)
					if ip == nil {
						return errors.New(v + " is not an IP")
					}
					ips = append(ips, ip)
				}
				cert.SetIPAddresses(ips)
				cert.SetDNSNames(c.StringSlice("domain"))
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
