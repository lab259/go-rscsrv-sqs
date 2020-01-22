COVERDIR=$(CURDIR)/.cover
COVERAGEFILE=$(COVERDIR)/cover.out

test:
	@go run github.com/onsi/ginkgo/ginkgo -race -failFast ./...

test-watch:
	@ginkgo watch -cover -r ./...

coverage-ci:
	@mkdir -p $(COVERDIR)
	@ginkgo --failFast -r -covermode=count --cover --trace ./
	@echo "mode: count" > "${COVERAGEFILE}"
	@find . -type f -name *.coverprofile -exec grep -h -v "^mode:" {} >> "${COVERAGEFILE}" \; -exec rm -f {} \;

coverage: coverage-ci
	@sed -i -e "s|_$(CURDIR)/|./|g" "${COVERAGEFILE}"

coverage-html:
	@go tool cover -html="${COVERAGEFILE}" -o .cover/report.html

fmt:
	@go fmt ./...

vet:
	@go vet ./...

dcup:
	@docker-compose up -d

dcdn:
	@docker-compose down

.PHONY: test test-watch coverage coverage-ci coverage-html fmt vet dcup dcdn
