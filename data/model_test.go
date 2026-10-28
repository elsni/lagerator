package data

import "testing"

// TestDataTableAcceptsPlainPayload verifies that storage does not require rendering methods.
func TestDataTableAcceptsPlainPayload(t *testing.T) {
	type payload struct {
		Value string `json:"value"`
	}
	table := NewDataTable[payload]()
	table.Add(Dataset[payload]{ID: 1, Name: "Plain", Data: payload{Value: "stored"}})
	set, found := table.GetPtr(1)
	if !found || set.Data.Value != "stored" {
		t.Fatalf("plain payload lookup failed: %+v, found=%v", set, found)
	}
}
