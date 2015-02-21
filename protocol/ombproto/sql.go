package ombproto

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func GetCreateSql() (string, error) {
	// Returns the SQL command that is used to create the pubrecord.db
	// We figure out where that file is by using GOPATH

	route := "src/github.com/soapboxsys/ombudslib/protocol/create_table.sql"

	gopath := os.Getenv("GOPATH")
	if len(gopath) < 1 {
		return "", errors.New("The EnvVar $GOPATH is not set!!")
	}

	fpath := filepath.Join(gopath, route)

	file, err := os.Open(fpath)
	if err != nil {
		return "", nil
	}

	sql := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		sql += "\n" + scanner.Text()
	}
	if len(sql) < 10 {
		return "", fmt.Errorf("File is empty")
	}
	return sql, nil
}
