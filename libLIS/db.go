package libLIS

import "database/sql"

var DB *sql.DB

func lastInsertID(tx *sql.Tx) (int, error) {
	id := 0
	sql_line := "SELECT LAST_INSERT_ID();"
	err := tx.QueryRow(sql_line).Scan(&id)
	if err != nil {
		TheLogger.Error(err)
		return 0, err
	}
	return id, nil
}

func lastInsertIDIfNecessary(tx *sql.Tx, id *int) error {
	if id == nil || *id != 0 {
		return nil
	}
	new_id, err := lastInsertID(tx)
	if err == nil {
		*id = new_id
	}
	return nil
}
