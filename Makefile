TARGET  := ui
  
GOBUILD := go build
  
everything : $(TARGET)
  
all : $(TARGET)
  
.PHONY: ui_db
ui_db :
	docker build -t ui_test_pg ./pg
	docker run --name ui_test_db -e POSTGRES_PASSWORD=ui_test -p 5432:5432 -d ui_test_pg
	
.PHONY: table_users
table_users:
	docker exec -td ui_test_db mkdir /opt/pg
	docker cp ./pg/users.sql ui_test_db:/opt/pg/
	docker exec -td ui_test_db bash -c 'psql -U ui_test < /opt/pg/users.sql'

.PHONY: clean         
clean :
	rm -rf ./output
	docker stop ui_test_db | true
	docker rm ui_test_db | true