run:
	go run cmd/restserver/main.go

test:
	go test -cover -race ./...

compose-up:
	docker compose up -d

compose-stop:
	docker compose stop

docker-exec:
	docker exec -it cashbycard /bin/bash

mockary:
	~/go/bin/mockery --all

create-volume:
	docker volume create --name=mysql_cashbycard_data

remove-volume:
	docker volume rm mysql_cashbycard_data
