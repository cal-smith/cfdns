package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/cloudflare/cloudflare-go"
	"golang.org/x/term"
)

// global flags
var token = flag.String("token", "", "Cloudflare API token. Must have edit access to all zone:domain pairs specified")
var skipHuman = flag.Bool("y", false, "Set this flag to skip confirmation of changes")

func getIP() string {
	res, err := http.Get("https://checkip.amazonaws.com/")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	ipBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	ip := strings.TrimSpace(string(ipBytes[:]))
	return ip
}

func getRecordsForZone(api *cloudflare.API, ctx context.Context, zoneName string) []cloudflare.DNSRecord {
	zoneId, err := api.ZoneIDByName(zoneName)

	if err != nil {
		log.Fatal(err)
	}

	records, err := api.DNSRecords(ctx, zoneId, cloudflare.DNSRecord{})

	if err != nil {
		log.Fatal(err)
	}

	return records
}

func matchesAnyDomain(recordName string, domains []string) bool {
	for _, domain := range domains {
		if strings.HasPrefix(recordName, domain) {
			return true
		}
	}
	return false
}

func askAHumanFrom(readFrom io.Reader, question string) bool {
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		return false
	}

	reader := bufio.NewReader(readFrom)

	fmt.Printf("%s (Y/n) ", question)
	res, _ := reader.ReadString('\n')
	res = strings.ToLower(strings.TrimSpace(res))
	if res == "y" || res == "yes" {
		return true
	} else {
		return false
	}
}

func askAHuman(question string) bool {
	return askAHumanFrom(os.Stdin, question)
}

func parseDomainConf(baseConfig string) map[string][]string {
	pairs := strings.Split(baseConfig, ",")
	zoneMap := make(map[string][]string)

	for _, pair := range pairs {
		parsed := strings.Split(pair, ":")
		zone := parsed[0]
		domain := parsed[1]
		zoneMap[zone] = append(zoneMap[zone], domain)
	}
	return zoneMap
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage for %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintln(flag.CommandLine.Output(), "  [zone:domain]")
		fmt.Fprintln(flag.CommandLine.Output(), "\tThe domains to update are provided as a comma separated list of zone:domian pairs")
		fmt.Fprintln(flag.CommandLine.Output(), "\tFor example: foo.com:cloud,bar.net:www")
	}
	flag.Parse()
	// config looks like - zone:domain,zone:domain
	config := flag.Arg(0)

	zoneMap := parseDomainConf(config)

	api, err := cloudflare.NewWithAPIToken(*token)

	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	currentIp := getIP()

	for zone, domains := range zoneMap {
		records := getRecordsForZone(api, ctx, zone)

		for _, record := range records {
			if matchesAnyDomain(record.Name, domains) {
				log.Printf("%s: %s %s last updated: %s\n", record.Name, record.Type, record.Content, record.ModifiedOn)
				log.Printf("detected %s, current %s", currentIp, record.Content)
				if currentIp == record.Content {
					log.Println("content unchanged. skipping update")
				} else if *skipHuman || askAHuman("continue") {
					err := api.UpdateDNSRecord(ctx, record.ZoneID, record.ID, cloudflare.DNSRecord{
						Type:    "A",
						Content: currentIp,
					})
					if err != nil {
						log.Fatal(err)
					}
					log.Println("done")
				}
			}
		}
	}
}
