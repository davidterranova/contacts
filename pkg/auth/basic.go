package auth

import (
	"crypto/sha1"
	"encoding/base64"
	"strings"

	"github.com/davidterranova/contacts/pkg/user"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func GrantAnyAccess() func(authToken string) (user.User, error) {
	return func(authToken string) (user.User, error) {
		reqUsername, _, ok := parseBasicAuth(authToken)
		if !ok {
			return user.Unauthenticated, ErrUnauthorized
		}

		id, err := uuid.FromBytes(sha1.New().Sum([]byte(reqUsername))[:16])
		log.Info().Str("username", reqUsername).Str("id", id.String()).Msg("granting access")
		return user.New(id), err
	}
}

func BasicAuth(username string, password string) func(authToken string) (user.User, error) {
	return func(authToken string) (user.User, error) {
		reqUsername, reqPassword, ok := parseBasicAuth(authToken)
		if !ok {
			return user.Unauthenticated, ErrUnauthorized
		}

		if reqUsername != username || reqPassword != password {
			return user.Unauthenticated, ErrUnauthorized
		}

		id := uuid.NewSHA1(uuid.NameSpaceOID, []byte(username))
		return user.New(id), nil
	}
}

func parseBasicAuth(auth string) (username string, password string, ok bool) {
	const prefix = "Basic "
	// Case insensitive prefix match. See Issue 22736.
	if len(auth) < len(prefix) || !equalFold(auth[:len(prefix)], prefix) {
		return "", "", false
	}
	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return "", "", false
	}
	cs := string(c)
	username, password, ok = strings.Cut(cs, ":")
	if !ok {
		return "", "", false
	}
	return username, password, true
}

func equalFold(s, t string) bool {
	if len(s) != len(t) {
		return false
	}
	for i := 0; i < len(s); i++ {
		if lower(s[i]) != lower(t[i]) {
			return false
		}
	}
	return true
}

func lower(b byte) byte {
	if 'A' <= b && b <= 'Z' {
		return b + ('a' - 'A')
	}
	return b
}
