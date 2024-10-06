package get

import (
	"database/sql"
	"footballresult/get/footballdata"
)

func Events(db *sql.DB) {

	footballdata.Start(db)

}
