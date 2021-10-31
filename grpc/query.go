package grpc

import (
	"context"
	"database/sql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/yavosh/sqln/proto"
)

func convert(input []string) []interface{} {
	res := make([]interface{}, len(input))
	for i := range input {
		res[i] = input[i]
	}
	return res
}

// ExecuteQuery will run the query and substitute any parameters provided
func (s *Server) ExecuteQuery(ctx context.Context, req *pb.QueryRequest) (*pb.QueryResult, error) {
	rows, err := s.db.QueryContext(ctx, req.Query, convert(req.Params)...)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	defer rows.Close()

	cols, err := rows.ColumnTypes()

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	rowsResult := make([]*pb.Row, 0, 1)
	rowResult := make([][]byte, len(cols))
	dest := make([]interface{}, len(cols)) // A temporary interface{} slice
	for i := range rowResult {
		dest[i] = &rowResult[i] // Put pointers to each string in the interface slice
	}

	for rows.Next() {
		err := rows.Scan(dest...)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}

		rowsResult = append(rowsResult, convertRow(cols, rowResult))
	}

	return &pb.QueryResult{
		Columns: convertColumns(cols),
		Rows:    rowsResult,
	}, nil
}

func convertRow(cols []*sql.ColumnType, rows [][]byte) *pb.Row {
	row := &pb.Row{
		Values: make([]*pb.Value, len(cols)),
	}

	for i := range rows {
		row.Values[i] = &pb.Value{Payload: rows[i]}
	}

	return row
}

func convertColumns(cols []*sql.ColumnType) []*pb.Column {
	res := make([]*pb.Column, len(cols))
	for i, c := range cols {
		res[i] = &pb.Column{
			Name: c.Name(),
			Type: 0, // TODO: types are in progress
		}
	}
	return res
}
