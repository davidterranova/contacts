package exports

import (
	"fmt"
	"html/template"
	"io"
	"strings"

	"github.com/davidterranova/contacts/internal/contacts/domain"
)

type BasicContactRenderer struct{}

var basicContactTemplate = template.Must(template.New("basicContact").Parse(`
<html>
  <head>
    <title>{{ .FirstName }} {{ .LastName }}</title>
  </head>
  <body>
    <h1>{{ .FirstName }} {{ .LastName }}</h1>
    <h2>{{ .Email }}</h2>
    <h2>{{ .Phone }}</h2>
    <p>
      {{ .CreatedAt }}
    </p>
    <p>
      {{ .UpdatedAt  }}
    </p>
  </body>
</html>
`))

func (BasicContactRenderer) Render(contact *domain.Contact) (io.Writer, error) {
	writer := new(strings.Builder)
	err := basicContactTemplate.Execute(writer, contact)
	if err != nil {
		return nil, fmt.Errorf("failed to render template: %w", err)
	}

	return writer, nil
}
