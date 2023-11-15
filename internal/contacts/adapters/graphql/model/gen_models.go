// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Contact struct {
	ID               string `json:"id"`
	CreatedAt        string `json:"createdAt"`
	UpdatedAt        string `json:"updatedAt"`
	FirstName        string `json:"firstName"`
	LastName         string `json:"lastName"`
	Phone            string `json:"phone"`
	Email            string `json:"email"`
	AggregateVersion int    `json:"aggregateVersion"`
}

type NewContact struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}