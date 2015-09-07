package mock

import (
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/mitch000001/go-harvest/harvest"
)

func TestNewUserService(t *testing.T) {
	mockUserService := UserService{
		Users: []*harvest.User{
			&harvest.User{ID: 1, UpdatedAt: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)},
		},
	}

	userService := NewUserService(mockUserService)

	if userService == nil {
		t.Logf("Expected userService not to be nil\n")
		t.FailNow()
	}

	expectedUsers := []*harvest.User{
		&harvest.User{ID: 1, UpdatedAt: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)},
	}
	var actualUsers []*harvest.User

	err := userService.All(&actualUsers, nil)

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	if !reflect.DeepEqual(expectedUsers, actualUsers) {
		t.Logf("Expected users to equal\n%q\n\tgot\n%q\n", expectedUsers, actualUsers)
		t.Fail()
	}
}

func TestUserServiceAll(t *testing.T) {
	mockUserService := UserService{
		Users: []*harvest.User{
			&harvest.User{ID: 1, UpdatedAt: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)},
			&harvest.User{ID: 2, UpdatedAt: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)},
		},
	}

	expectedUsers := []*harvest.User{
		&harvest.User{ID: 1, UpdatedAt: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)},
		&harvest.User{ID: 2, UpdatedAt: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)},
	}

	var actualUsers []*harvest.User

	err := mockUserService.All(&actualUsers, nil)

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	if !reflect.DeepEqual(expectedUsers, actualUsers) {
		t.Logf("Expected users to equal\n%q\n\tgot\n%q\n", expectedUsers, actualUsers)
		t.Fail()
	}
}

func TestUserServiceFind(t *testing.T) {
	mockUserService := UserService{
		Users: []*harvest.User{
			&harvest.User{ID: 1, UpdatedAt: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)},
			&harvest.User{ID: 2, UpdatedAt: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)},
		},
	}

	expectedUser := &harvest.User{ID: 1, UpdatedAt: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)}

	var actualUser harvest.User

	err := mockUserService.Find(1, &actualUser, nil)

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	if !reflect.DeepEqual(expectedUser, &actualUser) {
		t.Logf("Expected user to equal\n%#v\n\tgot\n%#v\n", expectedUser, &actualUser)
		t.Fail()
	}
}

func TestUserServiceCreate(t *testing.T) {
	mockUserService := UserService{
		Users: []*harvest.User{
			&harvest.User{ID: 1, UpdatedAt: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)},
		},
	}

	expectedUsers := []*harvest.User{
		&harvest.User{ID: 1, UpdatedAt: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)},
		&harvest.User{ID: 2, UpdatedAt: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)},
	}

	userToCreate := &harvest.User{ID: 2, UpdatedAt: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)}

	err := mockUserService.Create(userToCreate)

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	if !reflect.DeepEqual(expectedUsers, mockUserService.Users) {
		t.Logf("Expected users to equal\n%q\n\tgot\n%q\n", expectedUsers, mockUserService.Users)
		t.Fail()
	}
}

func TestUserServiceUpdate(t *testing.T) {
	mockUserService := UserService{
		Users: []*harvest.User{
			&harvest.User{ID: 1, FirstName: "Max", UpdatedAt: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)},
			&harvest.User{ID: 2, FirstName: "Charlie", UpdatedAt: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)},
		},
	}

	userToUpdate := &harvest.User{ID: 1, FirstName: "Kevin", UpdatedAt: time.Date(2015, 1, 2, 0, 0, 0, 0, time.UTC)}

	err := mockUserService.Update(userToUpdate)

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	if mockUserService.Users[0].FirstName != userToUpdate.FirstName {
		t.Logf("Expected user1.Firstname to equal %q, got %q\n", userToUpdate.FirstName, mockUserService.Users[0].FirstName)
		t.Fail()
	}
	if mockUserService.Users[0].UpdatedAt.Before(time.Now().In(time.UTC).Add(-2 * time.Second)) {
		t.Logf("Expected user1.UpdatedAt to be modified, was not: %v\n", mockUserService.Users[0].UpdatedAt)
		t.Fail()
	}
	if mockUserService.Users[1].FirstName != "Charlie" {
		t.Logf("Expected user2.Firstname to equal 'Charlie', got %q\n", userToUpdate.FirstName)
		t.Fail()
	}
}

func TestUserServiceDelete(t *testing.T) {
	mockUserService := UserService{
		Users: []*harvest.User{
			&harvest.User{ID: 1, UpdatedAt: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)},
			&harvest.User{ID: 2, UpdatedAt: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)},
		},
	}

	expectedUsers := []*harvest.User{
		&harvest.User{ID: 1, UpdatedAt: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)},
	}

	userToDelete := &harvest.User{ID: 2, UpdatedAt: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)}

	err := mockUserService.Delete(userToDelete)

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	if !reflect.DeepEqual(expectedUsers, mockUserService.Users) {
		t.Logf("Expected users to equal\n%q\n\tgot\n%q\n", expectedUsers, mockUserService.Users)
		t.Fail()
	}
}

func TestUserServiceToggle(t *testing.T) {
	mockUserService := UserService{
		Users: []*harvest.User{
			&harvest.User{ID: 1, IsActive: true, UpdatedAt: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)},
			&harvest.User{ID: 2, IsActive: true, UpdatedAt: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)},
		},
	}

	userToToggle := &harvest.User{ID: 1, UpdatedAt: time.Date(2015, 1, 2, 0, 0, 0, 0, time.UTC)}

	err := mockUserService.Toggle(userToToggle)

	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	if mockUserService.Users[0].IsActive != userToToggle.IsActive {
		t.Logf("Expected user1.IsActive to be %t, got %t\n", userToToggle.IsActive, mockUserService.Users[0].IsActive)
		t.Fail()
	}
	if mockUserService.Users[0].UpdatedAt.Before(time.Now().In(time.UTC).Add(-2 * time.Second)) {
		t.Logf("Expected user1.UpdatedAt to be modified, was not: %v\n", mockUserService.Users[0].UpdatedAt)
		t.Fail()
	}
	if mockUserService.Users[1].IsActive != true {
		t.Logf("Expected user2.IsActive to be true, got false\n")
		t.Fail()
	}
}

func TestUserServicePath(t *testing.T) {
	mockUserService := UserService{}

	path := mockUserService.Path()

	if path != "users" {
		t.Logf("Expected path to return 'users', got %q\n", path)
		t.Fail()
	}
}

func TestUserServiceURL(t *testing.T) {
	mockUserService := UserService{}

	actualUrl := mockUserService.URL()

	expectedUrl := url.URL{}

	if expectedUrl.String() != actualUrl.String() {
		t.Logf("Expected URL to return %q, got %q\n", expectedUrl, actualUrl)
		t.Fail()
	}
}
