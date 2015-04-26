package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRetrieveRecords(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{
    "domain": "testdomain.ru",
    "records": [
        {
            "record_id": 1,
            "type": "A",
            "domain": "testdomain.ru",
            "fqdn": "testdomain.ru",
            "ttl": 20,
            "subdomain": "test",
            "content": "192.168.1.1",
            "priority": 500
        },
        {
            "record_id": 2,
            "type": "A",
            "domain": "testdomain.ru",
            "fqdn": "testdomain.ru",
            "ttl": 20,
            "subdomain": "test",
            "content": "192.168.1.1",
            "priority": ""
        }
    ],
    "success": "true"
}`)
	}))
	defer ts.Close()

	records, err := RetrieveDomainRecords(ts.URL, "TEST_PDD_TOKEN", "testdomain.ru")
	if err != nil {
		t.Error("Expected to parse response without errors")
		panic(err)
	}

	assert.Equal(t, len(records), 2, "Exactly two records are expected.")

	r := records[0]

	assert.Equal(t, r.Record_id, 1)
	assert.Equal(t, r.Domain, "testdomain.ru")
	assert.Equal(t, r.Type, "A")
	assert.Equal(t, r.TTL, 20)
	assert.Equal(t, r.FQDN, "testdomain.ru")
	assert.Equal(t, r.Subdomain, "test")
	assert.Equal(t, r.Content, "192.168.1.1")
	assert.Equal(t, 500, r.Priority)

	assert.Equal(t, records[1].Priority, 0)
}

func TestFilterRecordsByType(t *testing.T) {
	records := []DNSRecord{
		DNSRecord{1, "A", "", "test.com", "www.test.com", 0, 100, "www"},
		DNSRecord{2, "MX", "", "test.com", "www.test.com", 0, 100, "www"},
		DNSRecord{3, "TXT", "", "test.com", "www.test.com", 0, 100, "www"},
	}

	filtered := FilterRecordsByType(records, []string{"MX", "A"})
	assert.Equal(t, len(filtered), 2, "Only two records are expected.")
	assert.Equal(t, filtered[0].Type, "A")
	assert.Equal(t, filtered[1].Type, "MX")
}
