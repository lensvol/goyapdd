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

	if len(records) != 1 {
		t.Error("Only one record should be returned.")
	}

	r := records[0]

	assert.Equal(t, r.Record_id, 1)
	assert.Equal(t, r.Domain, "testdomain.ru")
	assert.Equal(t, r.Type, "A")
	assert.Equal(t, r.TTL, 20)
	assert.Equal(t, r.FQDN, "testdomain.ru")
	assert.Equal(t, r.Subdomain, "test")
	assert.Equal(t, r.Content, "192.168.1.1")
	assert.Equal(t, r.Priority, 500)
}
