package exports

import (
	"fmt"
	"io"
	"strings"

	"github.com/davidterranova/contacts/internal/contacts/domain"
	"github.com/emersion/go-vcard"
)

type VCardRenderer struct{}

func (VCardRenderer) Render(contact *domain.Contact) (io.Writer, error) {
	writer := new(strings.Builder)
	encoder := vcard.NewEncoder(writer)

	card := make(vcard.Card)

	card.SetValue(vcard.FieldName, contact.FirstName+" "+contact.LastName)
	card.SetValue(vcard.FieldFormattedName, contact.FirstName+" "+contact.LastName)
	card.SetValue(vcard.FieldEmail, contact.Email)
	card.SetValue(vcard.FieldTelephone, contact.Phone)

	vcard.ToV4(card)

	err := encoder.Encode(card)
	if err != nil {
		return nil, fmt.Errorf("failed to render vcard: %w", err)
	}

	return writer, nil
}
