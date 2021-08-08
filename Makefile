TARGET  := ui

OUTPUT := output

.PHONY:
all : build

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
	docker build -t ui_test_pg ./pg
	docker run --name ui_test_db -e POSTGRES_PASSWORD=ui_test -p 5432:5432 -d ui_test_pg
	
.PHONY: table_users
table_users :
	docker exec -td ui_test_db mkdir /opt/pg
	docker cp ./pg/users.sql ui_test_db:/opt/pg/
	docker exec -td ui_test_db bash -c 'psql -U ui_test < /opt/pg/users.sql'
	
.PHONY: ui_swagger
ui_swagger:
	docker run --name ui_swagger -p 9901:8080 --mount type=bind,source=$(pwd)/swagger,target=/app swaggerapi/swagger-ui

.PHONY: clean         
clean :
	rm -rf ./$(OUTPUT)
	docker stop ui_test_db | true
	docker rm ui_test_db | true
	docker stop ui_swagger | true
	docker rm ui_swagger | true