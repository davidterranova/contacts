package http

import "github.com/davidterranova/contacts/internal/domain"

type Contact struct {
	Id        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

func fromDomain(c *domain.Contact) *Contact {
	return &Contact{
		Id:        c.Id.String(),
		CreatedAt: c.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: c.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		FirstName: c.FirstName,
		LastName:  c.LastName,
		Email:     c.Email,
		Phone:     c.Phone,
	}
}

func fromDomainList(contacts []*domain.Contact) []*Contact {
	var list = make([]*Contact, 0, len(contacts))
	for _, c := range contacts {
		list = append(list, fromDomain(c))
	}

	return list
}
