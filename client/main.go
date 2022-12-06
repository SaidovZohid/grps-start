package main

import (
	"context"
	"fmt"

	pb "github.com/SaidovZohid/grpc-student-server/genproto/user_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	client := pb.NewStudentClient(conn)

	// student, err := client.CreateStudent(context.Background(), &pb.StudentReq{
	// 	FirstName: "Zohid",
	// 	LastName:  "Saidov",
	// 	Email:     "zohidsaidov17+1@gmail.com",
	// 	Password:  "zohid2004",
	// })
	// student, err := client.GetStudent(context.Background(), &pb.IdMsg{Id: 1})
	// student, err := client.UpdateStudent(context.Background(), &pb.StudentRes{
	// 	Id: 1,
	// 	FirstName: "Zufar",
	// 	LastName: "Saidov",
	// 	Email: "zufarsaidov2000@gmail.com",
	// 	Password: "zufar2000", 
	// })
	s, err := client.DeleteStudent(context.Background(), &pb.IdMsg{Id: 1})
	if err != nil {
		panic(err)
	}
	fmt.Println(s)
	students, err := client.GetAllStudents(context.Background(), &pb.GetAllParamsReq{
		Limit: 10,
		Page: 1,
	})
	if err != nil {
		panic(err)
	}
	for _, student := range students.Students {
		fmt.Println(student)
	}
}
