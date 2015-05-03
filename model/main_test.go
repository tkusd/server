package model

type fixtureUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var fixtureUsers = []fixtureUser{
	fixtureUser{Name: "John", Email: "john@abc.com", Password: "123456"},
	fixtureUser{Name: "Mary", Email: "mary@abc.com", Password: "234567"},
}
