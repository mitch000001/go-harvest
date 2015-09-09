package mock

import (
	"net/url"
	"time"

	"github.com/mitch000001/go-harvest/harvest"
)

func NewUserService(userService UserService) *harvest.UserService {
	var service *harvest.UserService
	service = harvest.NewUserService(&userService, &userService)
	return service
}

type UserService struct {
	Users           []*harvest.User
	DayEntryService DayEntryService
}

func (u *UserService) All(users interface{}, params url.Values) error {
	*(users.(*[]*harvest.User)) = u.Users
	return nil
}

func (u *UserService) Find(id interface{}, user interface{}, params url.Values) error {
	ID := id.(int)
	for _, u := range u.Users {
		if ID == u.ID {
			*(user.(*harvest.User)) = *u
			return nil
		}
	}
	return nil
}

func (u *UserService) Create(model harvest.CrudModel) error {
	user := model.(*harvest.User)
	u.Users = append(u.Users, user)
	return nil
}

func (u *UserService) Update(model harvest.CrudModel) error {
	for _, user := range u.Users {
		if model.Id() == user.ID {
			*user = *model.(*harvest.User)
			user.UpdatedAt = time.Now().In(time.UTC)
		}
	}
	return nil
}

func (u *UserService) Delete(model harvest.CrudModel) error {
	var users []*harvest.User
	for _, user := range u.Users {
		if model.Id() != user.ID {
			users = append(users, user)
		}
	}
	u.Users = users
	return nil
}

func (u *UserService) Toggle(model harvest.ActiveTogglerCrudModel) error {
	for _, user := range u.Users {
		if model.Id() == user.ID {
			user.ToggleActive()
			user.UpdatedAt = time.Now().In(time.UTC)
		}
	}
	return nil
}

func (u *UserService) Path() string {
	return "users"
}

func (u *UserService) URL() url.URL {
	return url.URL{}
}

func (u *UserService) CrudEndpoint(path string) harvest.CrudEndpoint {
	return u.DayEntryService
}
