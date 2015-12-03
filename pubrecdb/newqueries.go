package pubrecdb

// BlockCount returns the number of blocks stored in the DB
func (db *PublicRecord) BlockCount() (int, error) {
	var count int
	query := `SELECT count(hash) FROM blocks;`
	err := db.conn.QueryRow(query).Scan(&count)
	if err != nil {
		return -1, err
	}

	return count, nil
}
