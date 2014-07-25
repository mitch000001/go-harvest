package main

import (
	"testing"
	"time"
)

func TestFindAllUsersUpdatedSince(t *testing.T) {
	client := createClient(t)
	updatedSince := time.Now().Add(-2 * time.Second)
	t.Logf("UpdatedSince: %+#v\n", updatedSince)
	users, err := client.Users.AllUpdatedSince(updatedSince)
	if err != nil {
		t.Fatalf("Got error %T with message: %s\n", err, err.Error())
	}
	for _, u := range users {
		t.Logf("User: '%+#v'\n", u)
	}
	if len(users) != 1 {
		t.Fatalf("Expected 1 user, got %d", len(users))
	}
}

func TestFindAllUsers(t *testing.T) {
	client := createClient(t)
	testAllFunc(client.Users.All, nil, t)
}

func TestFindUser(t *testing.T) {
	client := createClient(t)
	users, err := client.Users.All(nil)
	if err != nil {
		t.Fatalf("Got error %T with message: %s\n", err, err.Error())
	}
	first := users[0]
	testFindFunc(client.Users.Find, first.Id, t)
}

func TestCreateAndDeleteUser(t *testing.T) {
	client := createClient(t)
	user := User{
		FirstName: "Foo",
		LastName:  "Bar",
		Email:     "foo@example.com"}
	testCreateAndDeleteFunc(client.Users.Create, client.Users.Delete, &user, t)
}

func TestUpdateUser(t *testing.T) {
	client := createClient(t)
	user := &User{
		FirstName:  "Foo",
		LastName:   "Bar",
		Email:      "foo@example.com",
		Department: "Old Department",
	}
	user, err := client.Users.Create(user)
	if err != nil {
		panic(err)
	}
	defer client.Users.Delete(user)
	user.Department = "New Department"
	testUpdateFunc(client.Users.Update, user, "Department", t)

	// Wrong updates
	user.Email = "johjoo8 "
	testUpdateFuncInvalidUpdate(client.Users.Update, user, t)
}
