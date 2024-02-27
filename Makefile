#
# PostgreSQL Environment Variables
#
.EXPORT_ALL_VARIABLES:
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_USER ?= jonathan
DB_PASSWORD ?= john0804
DB_NAME ?= look_around
DATABASE_URL ?= sslmode=disable host=${DB_HOST} port=${DB_PORT} user=${DB_USER} password=${DB_PASSWORD} dbname=${DB_NAME}

#
# Ticket master API key
#
TICKET_MASTER_API_KEY ?= *****************

#
# Google Map API key
#
GOOGLE_MAP_API_KEY ?= ******************

#
# postgres
#
stop-pg:
	@echo "stop postgres..."
	@docker stop look-around-pg | true

restart-pg: stop-pg
	@echo "restart postgres..."
	@docker run -d --rm --name look-around-pg \
				-p 5432:5432 -e POSTGRES_DB=look_around \
				-e POSTGRES_USER=jonathan -e POSTGRES_PASSWORD=john0804 \
				postgres:13.4-alpine

#
# Assume we are all mac users
# Before using the Go bindings, you must install the libpostal C library. Make sure you have the following prerequisites:
#
install-prerequisites:
	@brew install curl autoconf automake libtool pkg-config || true
run:
	@echo "install ... "
	@go run main.go