package grpc_test

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/yavosh/sqln/grpc"
	"testing"

	pb "github.com/yavosh/sqln/proto"
)

func unwrap(v int64, err error) string {
	if err != nil {
		return fmt.Sprintf("err: %s", err)
	}
	return fmt.Sprintf("%d", v)
}

func TestServerExecuteQuery(t *testing.T) {

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf(err.Error(), err)
	}

	initialQueries := []string{
		"CREATE TABLE tab (id INTEGER PRIMARY KEY ASC, name VARCHAR(255), dob DATE);",
		"INSERT INTO tab(id, name, dob) VALUES(1, 'Alice', '2000-10-25');",
		"INSERT INTO tab(id, name, dob) VALUES(2, 'Jack', '2001-10-25');",
		"INSERT INTO tab(id, name, dob) VALUES(3, 'Jill', '2002-10-25');",
		"INSERT INTO tab(id, name, dob) VALUES(4, 'Dill', '2003-10-25');",
	}

	for _, q := range initialQueries {
		res, err := db.Exec(q)
		if err != nil {
			t.Fatalf(err.Error(), err)
		}
		fmt.Printf("%s %s\n", unwrap(res.LastInsertId()), unwrap(res.RowsAffected()))
	}

	server := grpc.NewServer(0, db)
	req := pb.QueryRequest{
		Db:     "main",
		Query:  "SELECT * FROM tab WHERE name = ?",
		Params: []string{"Alice"},
	}
	res, err := server.ExecuteQuery(context.Background(), &req)
	if err != nil {
		t.Fatalf(err.Error(), err)
	}

	fmt.Println(res)
}
