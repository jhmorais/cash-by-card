run:
	go run cmd/restserver/main.go

test:
	go test -cover -race ./...

compose-up:
	docker-compose up -d

compose-stop:
	docker-compose stop

docker-exec:
	docker exec -it devicemanager /bin/bash

mockary:
	~/go/bin/mockery --all

create-volume:
	docker volume create --name=mysql_devicemanager_data

remove-volume:
	docker volume rm mysql_devicemanager_data
