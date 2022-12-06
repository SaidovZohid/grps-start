package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/SaidovZohid/grpc-student-server/genproto/user_service"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	pb.UnimplementedStudentServer
	db *sqlx.DB
}

var (
	Host     string = "localhost"
	Port     string = "5432"
	User     string = "postgres"
	Password string = "1234"
	Database string = "students_db"
)

func main() {
	psqlUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
	Host,
	Port,
	User,
	Password,
	Database,
)
	psqlConn, err := sqlx.Connect("postgres", psqlUrl)

	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	fmt.Println(psqlUrl)
	fmt.Println("Connected Succesfully!")

	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		panic(err)
	}

	srv := grpc.NewServer()
	pb.RegisterStudentServer(srv, &Server{db: psqlConn})
	reflection.Register(srv)

	fmt.Println("Server Started")

	if e := srv.Serve(listener); e != nil {
		panic(e)
	}
}

func (d *Server) CreateStudent(ctx context.Context, s *pb.StudentReq) (*pb.StudentRes, error) {
	var res pb.StudentRes
	query := `
		INSERT INTO students (
			first_name,
			last_name,
			email,
			password
		) VALUES  ($1, $2, $3, $4)
		RETURNING 
		id, 
		first_name,
		last_name,
		email,
		password
	`

	err := d.db.QueryRow(
		query,
		s.FirstName,
		s.LastName,
		s.Email,
		s.Password,
	).Scan(
		&res.Id,
		&res.FirstName,
		&res.LastName,
		&res.Email,
		&res.Password,
	)
	if err != nil {
		return &pb.StudentRes{}, err
	}

	return &res, nil
}

func (d *Server) GetStudent(ctx context.Context, s *pb.IdMsg) (*pb.StudentRes, error) {
	query := `
		SELECT 
			id,
			first_name,
			last_name,
			email,
			password
		FROM students WHERE id = $1		
	`
	var res pb.StudentRes
	err := d.db.QueryRow(
		query,
		s.Id,
	).Scan(
		&res.Id,
		&res.FirstName,
		&res.LastName,
		&res.Email,
		&res.Password,
	)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (d *Server) UpdateStudent(ctx context.Context, s *pb.StudentRes) (*pb.StudentRes, error) {
	query := `
		UPDATE students SET 
			first_name = $1,
			last_name = $2,
			email = $3,
			password = $4
		WHERE id = $5		
	`
	_, err := d.db.Exec(
		query,
		s.FirstName,
		s.LastName,
		s.Email,
		s.Password,
		s.Id,
	)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (d *Server) DeleteStudent(ctx context.Context, s *pb.IdMsg) (*pb.Empty, error) {
	query := `
		DELETE FROM students WHERE id = $1
	`
	_, err := d.db.Exec(
		query,
		s.Id,
	)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (d *Server) GetAllStudents(ctx context.Context, params *pb.GetAllParamsReq) (*pb.GetAllStudentsRes, error) {
	offset := (params.Page - 1) * params.Limit
	limit := fmt.Sprintf(" LIMIT %d OFFSET %d ", params.Limit, offset)
	filter := ""
	if params.Search != "" {
		str := "%" + params.Search + "%"
		filter = fmt.Sprintf(" WHERE first_name ILIKE '%s' OR last_name ILIKE '%s' OR email ILIKE '%s'", str, str, str)
	}
	query := `
		SELECT 
			id,
			first_name,
			last_name,
			email,
			password
		FROM students 
	` + filter + limit
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}

	var res pb.GetAllStudentsRes
	res.Students = make([]*pb.StudentRes, 0)
	for rows.Next() {
		var s pb.StudentRes
		err := rows.Scan(
			&s.Id,
			&s.FirstName,
			&s.LastName,
			&s.Email,
			&s.Password,
		)
		if err != nil {
			return nil, err
		}
		res.Students = append(res.Students, &s)
	}

	queryCount := " SELECT count(1) FROM students " + filter
	err = d.db.QueryRow(queryCount).Scan(&res.Count)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
