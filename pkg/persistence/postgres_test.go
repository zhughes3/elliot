package persistence

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPostgresConnectionString(t *testing.T) {
	// newPostgresConnectionString(cfg dbConfig) string
	tcs := []struct {
		got      dbConfig
		expected string
	}{
		{
			got: dbConfig{
				User:     "boom",
				Password: "password",
				Host:     "somehost.com",
				Port:     "1234",
				Name:     "some_db_name",
			}, expected: "postgres://boom:password@somehost.com:1234/some_db_name?sslmode=disable",
		},
		{
			got: dbConfig{
				User:     "boom",
				Password: "password",
				Host:     "",
				Port:     "",
				Name:     "some_db_name",
			}, expected: "postgres://boom:password@:/some_db_name?sslmode=disable",
		},
		{
			got: dbConfig{
				User:     "",
				Password: "",
				Host:     "",
				Port:     "",
				Name:     "",
			}, expected: "postgres://:@:/?sslmode=disable",
		},
	}

	for _, tc := range tcs {
		assert.Equal(t, newPostgresConnectionString(tc.got), tc.expected, "identifiers should match")
	}
}
