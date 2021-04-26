package ch02

type Person struct {
	Id        int    `col:"id" json:"id" form:"id"`
	FirstName string `col:"first_name" json:"first_name" form:"first_name"`
	LastName  string `col:"last_name" json:"last_name" form:"last_name"`
}
