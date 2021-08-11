TARGET  := ui

OUTPUT := output

.PHONY:
all : build
	docker-compose up

.PHONY: build
build : $(OUTPUT)
	go build -o ./$(OUTPUT)/$(TARGET) ./main/
$(OUTPUT) :
	mkdir $(OUTPUT)

.PHONY: test
test :
	go test ./...
  
.PHONY: ui_db
ui_db :
	docker build -t ui_test_pg -f ./pg/Dockerfile .
	docker run --name ui_test_db -e POSTGRES_PASSWORD=ui_test -p 5432:5432 -d ui_test_pg
	
.PHONY: table_users
table_users :
	docker exec -td ui_test_db bash -c 'psql -U ui_test < /opt/pg/users.sql'
	
.PHONY: ui_swagger
ui_swagger :
	docker build -t ui_swagger_web -f ./swagger/Dockerfile .
	docker run -td --name ui_swagger -p 9901:8080 ui_swagger_web

.PHONY: clean         
clean :
	rm -rf ./$(OUTPUT)
	docker stop ui ui_swagger ui_test_db | true
	docker rm ui ui_swagger ui_test_db | true
	docker-compose down --rmi all -v
