// +build feature

package harvest_test

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
	t.Logf("Account company: %+#v\n", account.Company)
	t.Logf("Account user: %+#v\n", account.User)
}
