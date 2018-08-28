package main

import (
	"fmt"

	"github.com/gobuffalo/packr"
	"github.com/jackc/pgx"
	jsoniter "github.com/json-iterator/go"
	"github.com/urfave/cli"

	"log"
	"os"
)

// PerformInit actually performs init operation
func PerformInit(db *pgx.Conn, fhirVersion string) error {
	var schemaStatements []string

	box := packr.NewBox("./data")
	schema, err := box.MustBytes(fmt.Sprintf("schema/fhirbase-%s.sql.json", fhirVersion))

	if err != nil {
		log.Fatalf("Cannot find FHIR schema '%s'", fhirVersion)
	}

	err = jsoniter.Unmarshal(schema, &schemaStatements)

	if err != nil {
		log.Fatalf("Cannot parse FHIR schema '%s': %v", fhirVersion, err)
	}

	for _, stmt := range schemaStatements {
		_, err = db.Exec(stmt)

		if err != nil {
			log.Printf("PG error: %v\nWhile executing statement:\n%s\n", err, stmt)
		}
	}

	return nil
}

// InitCommand loads FHIR schema into database
func InitCommand(c *cli.Context) error {
	var fhirVersion string

	if c.NArg() > 0 {
		fhirVersion = c.Args().Get(0)
	} else {
		log.Printf("You must provide a FHIR version for `fhirbase init` command.")
		os.Exit(1)
	}

	db := GetConnection(nil)
	PerformInit(db, fhirVersion)

	log.Printf("Database initialized with FHIR schema version '%s'", fhirVersion)

	return nil
}
