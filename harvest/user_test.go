package main

import (
	"testing"
	"time"
)

func TestFindAllUsersUpdatedSince(t *testing.T) {
	client := createClient(t)
	users, err := client.Users.AllUpdatedSince(time.Now())
	if err != nil {
		t.Fatalf("Got error %T with message: %s\n", err, err.Error())
	}
	if len(users) != 0 {
		t.Fatalf("Expected 0 user, got %d", len(users))
	}
}

func TestFindAllUsers(t *testing.T) {
	client := createClient(t)
	users, err := client.Users.All()
	if err != nil {
		t.Fatalf("Got error %T with message: %s\n", err, err.Error())
	}
	if len(users) != 1 {
		t.Fatalf("Expected 1 user, got %d", len(users))
	}
	if users[0] == nil {
		t.Fatal("Expected user not to be nil")
	}
	for _, u := range users {
		t.Logf("User: %+#v\n", u)
	}
}

func TestFindUser(t *testing.T) {
	client := createClient(t)

	// Find existing user
	users, err := client.Users.All()
	if err != nil {
		t.Fatalf("Got error %T with message: %s\n", err, err.Error())
	}
	first := users[0]
	user, err := client.Users.Find(first.Id)
	if err != nil {
		t.Fatalf("Got error %T with message: %s\n", err, err.Error())
	}
	if first.Id != user.Id {
		t.Fatalf("Expect to find user with id %d, got user %#v\n", first.Id, user)
	}

	// No user with that id
	user, err = client.Users.Find(1)
	if err != nil {
		expectedErrorMessage := "No user found with id 1"
		if err.Error() != expectedErrorMessage {
			t.Fatalf("Expected ResponseError with message '%s', got error %T with message: %s\n", expectedErrorMessage, err, err.Error())
		}
	}
	if user != nil {
		t.Fatalf("Expected user to be nil, got '%+#v'", user)
	}
}

func TestCreateAndDeleteUser(t *testing.T) {
	client := createClient(t)
	user := User{
		FirstName: "Foo",
		LastName:  "Bar",
		Email:     "foo@example.com"}
	createdUser, err := client.Users.Create(&user)
	if err != nil {
		t.Fatalf("Got error %T with message: %s\n", err, err.Error())
	}
	if createdUser == nil {
		t.Fatal("Expected user not to be nil")
	}
	t.Logf("Got returned user: %+#v\n", createdUser)
	deleted, err := client.Users.Delete(&user)
	if err != nil {
		t.Fatalf("Got error %T with message: %s\n", err, err.Error())
	}
	if !deleted {
		t.Fatalf("Could not delete user created for test")
	}
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
	updatedUser, err := client.Users.Update(user)
	if err != nil {
		t.Fatalf("Got error %T with message: %s\n", err, err.Error())
	}
	if updatedUser.Department != user.Department {
		t.Fatalf("Expected updated department to equal '%s', got '%s'", user.Department, updatedUser.Department)
	}

	// Wrong updates
	user.FirstName = ""
	updatedUser, err = client.Users.Update(user)
	if err == nil {
		t.Fatal("Expected ResponseError, got nil")
	}
	if updatedUser != nil {
		t.Fatalf("Expected user to be nil, got '%+#v'", updatedUser)
	}

}
