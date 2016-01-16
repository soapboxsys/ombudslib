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

func (db *PublicRecord) CurrentTip() string {
	return "this is a fake chain tip string"
}
