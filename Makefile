CURRENT_DIR=${shell pwd}
DB_URL=postgresql://postgres:1234@localhost:5432/students_db?sslmode=disable

proto-gen:
	rm -rf genproto
	./gen-proto.sh ${CURRENT_DIR}

migrate-up:
	migrate -path migrations -database "$(DB_URL)" -verbose up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" -verbose down