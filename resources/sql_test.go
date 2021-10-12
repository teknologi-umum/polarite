package resources_test

import (
	"polarite/resources"
	"testing"
)

func TestParseURL_MySQL(t *testing.T) {
	u, err := resources.ParseURL("mysql://root:password@localhost:3306/polarite?tls=true")
	if err != nil {
		t.Error("an error was thrown:", err)
	}

	if u != "root:password@tcp(localhost:3306)/polarite?tls=true" {
		t.Error("u is not \"root:password@tcp(localhost:3306)/polarite?tls=true\", got:", u)
	}
}

func TestParseURL_Postgres(t *testing.T) {
	u, err := resources.ParseURL("postgres://root:password@localhost:5432/polarite?ssl=verify")
	if err != nil {
		t.Error("an error was thrown:", err)
	}

	if u != "user=root password=password host=localhost port=5432 dbname=polarite ssl=verify" {
		t.Error("u is not \"user=root password=password host=localhost port=5432 dbname=polarite ssl=verify\", got:", u)
	}
}

func TestParseURL_Other(t *testing.T) {
	u, err := resources.ParseURL("redis://user:password@localhost:6739/")
	if err != nil {
		t.Error("an error was thrown:", err)
	}

	if u != "redis://user:password@localhost:6739/" {
		t.Error("u is not \"redis://user:password@localhost:6739/\", got:", u)
	}
}
