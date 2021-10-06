package snowflake

import (
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func Organization(name string) *Builder {
	return &Builder{
		name:       name,
		entityType: OrganizationType,
	}
}

// organization is a go representation of a grant that can be used in conjunction
// with github.com/jmoiron/sqlx
type organization struct {
	Name                           string    `db:"name"`
	State                          string    `db:"state"`
	AdminName                      string    `db:"admin_name"`
	AdminPassword                  string    `db:"admin_password"`
	Email                          string    `db:"email"`
	Edition                        string    `db:"edition"`
	Region                         string    `db:"region"`
	FirstName                      string    `db:"first_name"`
	LastName                       string    `db:"last_name"`
	Comment                        string    `db:"comment"`
}

func ScanOrganization(row *sqlx.Row) (*organization, error) {
	w := &organization{}
	err := row.StructScan(w)
	return w, err
}

func ListOrganizations(db *sql.DB) ([]organization, error) {
	stmt := "SHOW ORGANIZATION ACCOUNTS"
	rows, err := Query(db, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dbs := []organization{}
	err = sqlx.StructScan(rows, &dbs)
	if err == sql.ErrNoRows {
		log.Printf("[DEBUG] no organizations found")
		return nil, nil
	}
	return dbs, errors.Wrapf(err, "unable to scan row for %s", stmt)
}
