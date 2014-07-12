package main

import (
	"testing"
)

func TestAccountInformation(t *testing.T) {
	client := createClient(t)
	account, err := client.Account()
	if err != nil {
		t.Fatalf("Got error %T with message: %s\n", err, err.Error())
	}
	if account == nil {
		t.Fatal("Expected account not to be nil")
	}
	t.Logf("Account: %+#v\n", account)
}
