BINDIR=${CURDIR}/bin
PACKAGE=${CURDIR}/cmd/app/
DSN=postgres://postgres:postgres@localhost:5432

.PHONY: migration-create migration-reset migration-up bindir build clean test up down


migration-reset:
	goose -dir "migrations/" postgres ${DSN} reset

migration-up:
	goose -dir "migrations/" postgres postgres://postgres:postgres@localhost:5432 up


bindir:
	mkdir -p ${BINDIR}

build: bindir
	go build -o ${BINDIR}/app ${PACKAGE}

run: build
	go run ./${BINDIR}

clean:
	rm -rf ${BINDIR}

test:
	go test  -v ./e2e -run TestBasicScenario

up:
	docker-compose up -d

down:
	docker-compose down

start:
	make test
	make build
	make up
