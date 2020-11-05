GOCMD=$(if $(shell which richgo),richgo,go)

test:
	$(GOCMD) test -cover ./...

test-watch:
	reflex -s --decoration=none -r \.go$$ -- make test
