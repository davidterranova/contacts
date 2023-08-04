package main

import (
	"github.com/rs/zerolog/log"

	"github.com/davidterranova/contacts/cmd"
)

// "github.com/rs/zerolog/log"

// "github.com/davidterranova/contacts/cmd"

func main() {
	err := cmd.Execute()
	if err != nil {
		log.
			Fatal().
			Err(err).
			Msg("failed to start contacts")
	}

	// store := eventsourcing.NewEventStore()
	// registry := registry.NewRegistry()
	// registry.Register("contact", domain.Contact{})
	// router := eventsourcing.NewRouter(registry, store)

	// aggregate, err := router.Route(eventsourcing.NewEventBase(
	// 	domain.NewContact(),
	// 	"contact_created",
	// 	domain.ContactCreated{},
	// ))
	// if err != nil {
	// 	panic(err)
	// }
	// print("contact_created", aggregate)

	// aggregate, err = router.Route(eventsourcing.NewEventBase(
	// 	aggregate,
	// 	"contact_email_updated",
	// 	domain.ContactEmailUpdated{
	// 		Email: "dterranova@contact.local",
	// 	},
	// ))
	// if err != nil {
	// 	panic(err)
	// }
	// print("contact_email_updated", aggregate)

	// aggregate, err = router.Route(eventsourcing.NewEventBase(
	// 	aggregate,
	// 	"contact_updated_name",
	// 	domain.ContactNameUpdated{
	// 		FirstName: "David",
	// 		LastName:  "Terranova",
	// 	},
	// ))
	// if err != nil {
	// 	panic(err)
	// }
	// print("contact_updated_name", aggregate)

	// contactNameUpdated := eventsourcing.NewEventBase(
	// 	aggregate,
	// 	"contact_updated_name",
	// 	domain.ContactNameUpdated{
	// 		FirstName: "Plop",
	// 		LastName:  "Plip",
	// 	},
	// )
	// aggregate, err = router.Route(contactNameUpdated)
	// if err != nil {
	// 	panic(err)
	// }
	// print("contact_updated_name", aggregate)

	// fmt.Println("----events for aggregate ", aggregate.AggregateId())
	// for _, e := range store.Get(aggregate.AggregateId()) {
	// 	fmt.Printf("----%s: %s\n", e.EventType(), e.Data())
	// }

	// cqrs.Test()
}

// func print(event string, a eventsourcing.Aggregate) {
// 	fmt.Printf("%s: type<%T>: %s\n\n", event, a, a)
// }
