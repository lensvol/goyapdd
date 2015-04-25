package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

type PriorityValue int64

type DNSRecord struct {
	Content   string
	Domain    string
	FQDN      string
	Priority  PriorityValue
	TTL       int64
	Subdomain string
	Record_id int64
	Type      string
}

type ListDNSRecordsResponse struct {
	Records []DNSRecord
	Success string
}

func (p *PriorityValue) UnmarshalJSON(b []byte) (err error) {
	s, n := "foobar", uint64(0)
	if err = json.Unmarshal(b, &s); err == nil {
		*p = 0
		return nil
	}

	if err = json.Unmarshal(b, &n); err == nil {
		_ = "breakpoint"
		*p = PriorityValue(n)
	}
	return nil
}

func PrintRecords(records []DNSRecord) {
	fmt.Printf("Type\t\tSubdomain\tContent\n")
	fmt.Printf("-----\t\t--------\t-------\n")

	for _, record := range records {
		fmt.Printf(
			"%-12s\t%-12s\t%-12s\n",
			record.Type,
			record.Subdomain,
			record.Content,
		)
	}
}

func RetrieveDomainRecords(apiURL string, pddToken string, domain string) ([]DNSRecord, error) {
	client := &http.Client{}
	req, _ := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/api2/admin/dns/list", apiURL),
		nil)
	req.Header.Set("PddToken", pddToken)

	values := req.URL.Query()
	values.Add("domain", domain)
	req.URL.RawQuery = values.Encode()

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var container ListDNSRecordsResponse
	err = json.Unmarshal(body, &container)

	if err != nil {
		return nil, err
	}

	return container.Records, nil
}

func main() {
	pddTokenPtr := flag.String("pdd-token", "<auth token>", "PDD authenthication ticket.")
	domainPtr := flag.String("domain", "<domain>", "Domain name.")
	flag.Parse()

	dnsRecords, err := RetrieveDomainRecords("https://pddimp.yandex.ru", *pddTokenPtr, *domainPtr)
	if err != nil {
		panic(err)
	}

	PrintRecords(dnsRecords)
}
