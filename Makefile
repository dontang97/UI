TARGET  := ui
  
GOBUILD := go build
  
everything : $(TARGET)
  
all : $(TARGET)
  
.PHONY: ui_db
ui_db :
	docker build -t ui_test_pg ./pg
	docker run --name ui_test_db -e POSTGRES_PASSWORD=ui_test -p 5432:5432 -d ui_test_pg

.PHONY: clean         
clean :
	rm -rf ./output
	docker stop ui_test_db | true
	docker rm ui_test_db | true