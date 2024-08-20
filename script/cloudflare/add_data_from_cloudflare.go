package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"github.com/48Club/ip-waf-helper/database"
	"github.com/48Club/ip-waf-helper/types"
	"github.com/cloudflare/cloudflare-go"
)

func main() {
	apiKey := os.Getenv("CLOUDFLARE_API_KEY")
	zone := os.Getenv("CLOUDFLARE_ZONE")
	api, err := cloudflare.NewWithAPIToken(apiKey)
	if err != nil {
		panic(err)
	}
	page := 1
	var rules = []cloudflare.AccessRule{}
	var ruleIds = []string{}
	tc := time.NewTicker((5 * time.Minute) / 1200)
	for {
		<-tc.C
		time.After(100 * time.Millisecond)
		rule, err := api.ListAccountAccessRules(context.Background(), zone, cloudflare.AccessRule{}, page)
		if err != nil {
			log.Printf("Error: %s", err)
			continue
		}
		rules = append(rules, rule.Result...)
		if rule.ResultInfo.Page == rule.ResultInfo.TotalPages {
			break
		}
		page++
	}
	for _, rule := range rules {
		ip := rule.Configuration.Value
		if rule.Configuration.Target == "ip_range" {
			_, cidr64, _ := net.ParseCIDR(ip)
			ip = cidr64.String()
		}
		line := types.IPWaf{
			IP: ip,
		}
		tx := database.Server.FirstOrCreate(&line, line)
		if tx.Error != nil {
			panic(tx.Error)
		}
		ruleIds = append(ruleIds, rule.ID)
	}
	log.Printf("Done with %d rules", len(rules))
	for _, v := range ruleIds {
		<-tc.C
		_, err := api.DeleteAccountAccessRule(context.Background(), zone, v)
		if err != nil {
			log.Printf("Error: %s", err)
		}
	}
}
