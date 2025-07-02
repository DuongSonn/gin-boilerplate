# Varibles
GOCMD := go
AIRCMD := air
ATLCMD := atlas

FILE  ?= ""

dev:
	$(AIRCMD) 

run:
	$(GOCMD) run main.go

migrate:
	$(ATLCMD) migrate diff --env gorm

migrate-sync:
	$(ATLCMD) migrate hash