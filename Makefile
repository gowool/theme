ALL_GO_MOD_DIRS := $(shell find . -type f -name 'go.mod' -exec dirname {} \; | sort)

test:
	set -e; for dir in $(ALL_GO_MOD_DIRS); do \
	  echo "go test in $${dir}"; \
	  (cd "$${dir}" && \
	    go test && \
	    go vet); \
	done

go_mod_tidy:
	set -e; for dir in $(ALL_GO_MOD_DIRS); do \
	  echo "go mod tidy in $${dir}"; \
	  (cd "$${dir}" && \
	    go get -u ./... && \
	    go mod tidy); \
	done

fmt:
	gofmt -w -s ./
	goimports -w  -local github.com/gowool/theme ./