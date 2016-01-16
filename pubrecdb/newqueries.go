package pubrecdb

import "fmt"

// BlockCount returns the number of blocks stored in the DB
func (db *PublicRecord) BlockCount() (int, error) {
	return db.countRows("blocks")
}

func (db *PublicRecord) BulletinCount() (int, error) {
	return db.countRows("bulletins")
}

func (db *PublicRecord) EndoCount() (int, error) {
	return db.countRows("endorsements")
}

func (db *PublicRecord) countRows(table string) (int, error) {
	var count int
	query := fmt.Sprintf(`SELECT count(*) FROM %s;`, table)
	err := db.conn.QueryRow(query).Scan(&count)
	if err != nil {
		return -1, err
	}

	return count, nil
}

func (db *PublicRecord) CurrentTip() (string, error) {
	query := "SELECT hash FROM blocks ORDER BY height DESC LIMIT 1"
	var hash string
	err := db.conn.QueryRow(query).Scan(&hash)
	if err != nil {
		return "", err
	}

	return hash, nil
}
