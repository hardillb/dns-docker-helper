package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/miekg/dns"
	containertypes "github.com/moby/moby/api/types/container"
	eventstypes "github.com/moby/moby/api/types/events"
	"github.com/moby/moby/client"
)

var (
	dnsClient dns.Client
	cname     string
	dnsServer string
	keyName   string
	keySecret string
	zone      string
)

func main() {
	dnsServer = os.Getenv("DNS_SERVER")
	cname = os.Getenv("CNAME_TARGET")
	keyName = os.Getenv("KEY_NAME")
	keySecret = os.Getenv("KEY_SECRET")
	zone = os.Getenv("ZONE")

	if dnsServer == "" ||
		cname == "" ||
		keyName == "" ||
		keySecret == "" ||
		zone == "" {
		panic("Env vars must be set")
	}

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	dnsClient := new(dns.Client)
	dnsClient.TsigSecret = map[string]string{dns.Fqdn(keyName): keySecret}

	containers, err := cli.ContainerList(ctx, containertypes.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Println(container.ID)
		details, err := cli.ContainerInspect(ctx, container.ID)
		if err != nil {
			panic(err)
		}
		for _, env := range details.Config.Env {
			if strings.HasPrefix(env, "VIRTUAL_HOST=") {
				host := strings.Split(env, "=")[1]
				fmt.Println(host)
				// need to check for , to see if multiple names
				// Also need to check for existing names
				if strings.Contains(host, ",") {
					hosts := strings.Split(host, ",")
					for _, element := range hosts {
						addRecord(element)
					}
				} else {
					addRecord(host)
				}
			}
		}
	}

	msgs, errs := cli.Events(ctx, eventstypes.ListOptions{})

	for {
		select {
		case err := <-errs:
			fmt.Println(err)
		case msg := <-msgs:
			if msg.Action == "start" {
				details, err := cli.ContainerInspect(ctx, msg.Actor.ID)
				if err != nil {
					panic(err)
				}
				for _, env := range details.Config.Env {
					if strings.HasPrefix(env, "VIRTUAL_HOST=") {
						host := strings.Split(env, "=")[1]
						fmt.Println(host)
						if strings.Contains(host, ",") {
							hosts := strings.Split(host, ",")
							for _, element := range hosts {
								addRecord(element)
							}
						} else {
							addRecord(host)
						}
					}
				}
			} else if msg.Action == "kill" || msg.Action == "stop" {
				details, err := cli.ContainerInspect(ctx, msg.Actor.ID)
				if err != nil {
					panic(err)
				}
				for _, env := range details.Config.Env {
					if strings.HasPrefix(env, "VIRTUAL_HOST=") {
						host := strings.Split(env, "=")[1]
						fmt.Println(host)
						if strings.Contains(host, ",") {
							hosts := strings.Split(host, ",")
							for _, element := range hosts {
								removeRecord(element)
							}
						} else {
							removeRecord(host)
						}
					}
				}
			}
		}
	}
}

func addRecord(name string) {
	update, err := dns.NewRR(fmt.Sprintf("%s 300 IN CNAME %s", name, cname))
	if err != nil {
		panic(err)
	}
	updates := make([]dns.RR, 1)
	updates[0] = update
	message := new(dns.Msg)
	message.SetUpdate(dns.Fqdn(zone))
	message.Insert(updates)
	message.SetTsig(dns.Fqdn(keyName), dns.HmacSHA256, 300, time.Now().Unix())
	in, rtt, err := dnsClient.Exchange(message, dnsServer)
	if err != nil {
		fmt.Printf("%v in %d\n", in, rtt)
	} else {
		fmt.Printf("%v\n", err)
	}

}

func removeRecord(name string) {
	update, err := dns.NewRR(fmt.Sprintf("%s 300 IN CNAME %s", name, cname))
	if err != nil {
		panic(err)
	}
	updates := make([]dns.RR, 1)
	updates[0] = update
	message := new(dns.Msg)
	message.SetUpdate(dns.Fqdn(zone))
	message.Remove(updates)
	message.SetTsig(dns.Fqdn(keyName), dns.HmacSHA256, 300, time.Now().Unix())
	in, rtt, err := dnsClient.Exchange(message, dnsServer)
	if err != nil {
		fmt.Printf("%v in %d\n", in, rtt)
	} else {
		fmt.Printf("%v\n", err)
	}
}

func checkRecordExists(name string) bool {
	return false
}
