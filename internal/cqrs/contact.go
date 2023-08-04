package cqrs

// import (
// 	"errors"
// 	"fmt"
// 	"time"

// 	"github.com/davidterranova/contacts/internal/cqrs/registry"
// 	"github.com/google/uuid"
// )

// const AggregateContact AggregateType = "contact"

// var (
// 	ErrAggregateAlreadyExists = errors.New("aggregate already exists")
// 	ErrAggregateNotFound      = errors.New("aggregate not found")
// 	ErrInvalidAggregateType   = errors.New("invalid aggregate type")
// )

// type Contact struct {
// 	Id        uuid.UUID
// 	CreatedAt time.Time

// 	FirstName string
// 	LastName  string
// 	Email     string
// }

// func (c Contact) AggregateId() uuid.UUID {
// 	return c.Id
// }

// func (c Contact) AggregateType() AggregateType {
// 	return AggregateContact
// }

// type CmdUpdateContactName struct {
// 	createdAt time.Time

// 	ContactId uuid.UUID

// 	FirstName string
// 	LastName  string
// }

// func (c CmdUpdateContactName) AggregateId() uuid.UUID {
// 	return c.ContactId
// }

// func (c CmdUpdateContactName) AggregateType() AggregateType {
// 	return AggregateContact
// }

// func (c CmdUpdateContactName) CreatedAt() time.Time {
// 	return c.createdAt
// }

// func (c CmdUpdateContactName) Apply(aggregate Aggregate) ([]Event[Aggregate], error) {
// 	if aggregate.AggregateId() == uuid.Nil {
// 		return nil, ErrAggregateNotFound
// 	}

// 	return []Event[Aggregate]{
// 		EvtContactNameUpdated{
// 			id:        uuid.New(),
// 			createdAt: c.createdAt,
// 			contactId: c.ContactId,
// 			firstName: c.FirstName,
// 			lastName:  c.LastName,
// 		},
// 	}, nil
// }

// type EvtContactNameUpdated struct {
// 	id        uuid.UUID
// 	contactId uuid.UUID
// 	createdAt time.Time

// 	firstName string
// 	lastName  string
// }

// func (e EvtContactNameUpdated) Id() uuid.UUID {
// 	return e.id
// }

// func (e EvtContactNameUpdated) AggregateId() uuid.UUID {
// 	return e.contactId
// }

// func (e EvtContactNameUpdated) EventType() string {
// 	return "contact.name-updated"
// }

// func (e EvtContactNameUpdated) CreatedAt() time.Time {
// 	return e.createdAt
// }

// func (e EvtContactNameUpdated) Apply(aggregate Aggregate) error {
// 	contact, ok := aggregate.(*Contact)
// 	if !ok {
// 		return ErrInvalidAggregateType
// 	}

// 	contact.Id = e.contactId
// 	contact.FirstName = e.firstName
// 	contact.LastName = e.lastName

// 	return nil
// }

// func Test() {
// 	registry := registry.NewRegistry[Aggregate]()
// 	err := registry.Register("contact", Contact{})
// 	if err != nil {
// 		panic(err)
// 	}

// 	eventStore := NewEventStore[Aggregate]()
// 	commandHandler := NewCommandHandler[Contact](eventStore, registry)

// 	createContact := CmdCreateContact{
// 		createdAt: time.Now(),
// 		ContactId: uuid.New(),
// 		FirstName: "John",
// 		LastName:  "Doe",
// 		Email:     "test@contacts.local",
// 	}

// 	contact1, err := commandHandler.Handle(createContact)
// 	if err != nil {
// 		panic(err)
// 	}

// 	createContact.ContactId = uuid.New()
// 	contact2, err := commandHandler.Handle(createContact)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Printf("Contact1: %+v\n", contact1)
// 	fmt.Printf("Contact2: %+v\n", contact2)

// 	updateContactName := CmdUpdateContactName{
// 		createdAt: time.Now(),
// 		ContactId: contact1.AggregateId(),
// 		FirstName: "Sam",
// 		LastName:  "Smith",
// 	}
// 	contact1, err = commandHandler.Handle(updateContactName)
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Printf("Contact1: %+v\n", contact1)
// }
