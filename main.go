package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

type DNSRecord struct {
	Content   string
	Domain    string
	FQDN      string
	Priority  int64
	TTL       int64
	Subdomain string
	Record_id int64
	Type      string
}

type DNSRecords struct {
	Records []DNSRecord
	Success string
}

func PrintRecords(records []DNSRecord) {
	for _, record := range records {
		fmt.Printf(
			"%s %s %s\n",
			record.Subdomain,
			record.Type,
			record.Content,
		)
	}
}

func main() {
	pddTokenPtr := flag.String("pdd-token", "<auth token>", "PDD authenthication ticket.")
	domainPtr := flag.String("domain", "<domain>", "Domain name.")
	flag.Parse()

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://pddimp.yandex.ru/api2/admin/dns/list", nil)
	req.Header.Set("PddToken", *pddTokenPtr)

	values := req.URL.Query()
	values.Add("domain", *domainPtr)
	req.URL.RawQuery = values.Encode()

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)

	var container DNSRecords
	err = json.Unmarshal(body, &container)
	PrintRecords(container.Records)

	defer resp.Body.Close()
}
