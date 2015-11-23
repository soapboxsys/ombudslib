package pubrecdb

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestSchemaMatch(t *testing.T) {
	schema_path := os.Getenv("GOPATH") + "/src/github.com/soapboxsys/ombudslib/pubrecdb/schema.sql"
	f, err := os.Open(filepath.Clean(schema_path))
	if err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	file_sql := string(b)
	func_sql := createSql()

	if file_sql != func_sql {
		t.Fatal(fmt.Errorf("schema.sql and createSql() do not match. The intended DB"+
			" schema has changed. \n====schema.sql====:\n%b\n====createSql()====:\n%b",
			[]byte(file_sql), []byte(func_sql),
		))
	}

}
