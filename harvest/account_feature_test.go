// +build feature

package harvest_test

import (
	"testing"
)

func TestAccountInformation(t *testing.T) {
	client := createClient(t)
	account, err := client.Account()
	if err != nil {
		t.Logf("Got error %T with message: %s\n", err, err.Error())
		t.Fail()
	}
	if account == nil {
		t.Logf("Expected account not to be nil")
		t.Fail()
	}
	t.Logf("Account: %+#v\n", account)
	t.Logf("Account company: %+#v\n", account.Company)
	t.Logf("Account user: %+#v\n", account.User)
}
