package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
)

type PriorityValue int64
type DNSRecords []DNSRecord

type DNSRecord struct {
	Record_id int64
	Type      string
	Content   string
	Domain    string
	FQDN      string
	Priority  PriorityValue
	TTL       int64
	Subdomain string
}

type ListDNSRecordsResponse struct {
	Records DNSRecords
	Success string
}

func (r DNSRecords) Len() int {
	return len(r)
}
func (r DNSRecords) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}
func (r DNSRecords) Less(i, j int) bool {
	return r[i].Record_id < r[j].Record_id
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
	fmt.Printf("ID\t\tType\t\tSubdomain\tContent\n")
	fmt.Printf("--------\t-----\t\t--------\t-------\n")

	for _, record := range records {
		fmt.Printf(
			"%-12d\t%-12s\t%-12s\t%-12s\n",
			record.Record_id,
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

func Contains(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func FilterRecordsByType(records []DNSRecord, types []string) []DNSRecord {
	var approved_records []DNSRecord

	for _, r := range records {
		if Contains(r.Type, types) {
			approved_records = append(approved_records, r)
		}
	}
	return approved_records
}

func main() {
	pddTokenPtr := flag.String("pdd-token", "<auth token>", "PDD authenthication ticket.")
	domainPtr := flag.String("domain", "<domain>", "Domain name.")
	recTypesPtr := flag.String("filter-types", "A,MX,NS,SRV,TXT", "Filter types.")
	flag.Parse()

	var allowed_types []string

	if recTypesPtr != nil && len(*recTypesPtr) > 0 {
		allowed_types = strings.Split(*recTypesPtr, ",")
	} else {
		allowed_types = []string{"A", "MX", "SRV", "NS", "TXT"}
	}

	dnsRecords, err := RetrieveDomainRecords("https://pddimp.yandex.ru", *pddTokenPtr, *domainPtr)
	if err != nil {
		panic(err)
	}

	sort.Sort(DNSRecords(dnsRecords))
	dnsRecords = FilterRecordsByType(dnsRecords, allowed_types)
	PrintRecords(dnsRecords)
}
