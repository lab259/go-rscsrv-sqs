
GOPATH=$(CURDIR)/.gopath
GOPATHCMD=GOPATH=$(GOPATH)

COVERDIR=$(CURDIR)/.cover
COVERAGEFILE=$(COVERDIR)/cover.out

DOCKERCOMPOSETEST := docker-compose -f docker-compose.yml

.PHONY: deps deps-ci coverage coverage-ci test test-watch coverage coverage-html

test:
	@${GOPATHCMD} ginkgo --failFast ./...

test-watch:
	@${GOPATHCMD} ginkgo watch -cover -r ./...

coverage-ci:
	@mkdir -p $(COVERDIR)
	@${GOPATHCMD} ginkgo --failFast -r -covermode=count --cover --trace ./
	@echo "mode: count" > "${COVERAGEFILE}"
	@find . -type f -name *.coverprofile -exec grep -h -v "^mode:" {} >> "${COVERAGEFILE}" \; -exec rm -f {} \;

coverage: coverage-ci
	@sed -i -e "s|_$(CURDIR)/|./|g" "${COVERAGEFILE}"

coverage-html:
	@$(GOPATHCMD) go tool cover -html="${COVERAGEFILE}" -o .cover/report.html

fmt:
	@$(GOPATHCMD) go fmt

deps:
	@mkdir -p ${GOPATH}
	@go list -f '{{join .Deps "\n"}}' . | xargs go list -f '{{if not .Standard}}{{.ImportPath}}{{end}}' | GOPATH=${GOPATH} xargs go get -v -t
	@go list -f '{{join .TestImports "\n"}}' . | xargs go list -f '{{if not .Standard}}{{.ImportPath}}{{end}}' | GOPATH=${GOPATH} xargs go get -v -t

deps-ci:
	-go get -v -t ./...

dco-test-up:
	${DOCKERCOMPOSETEST} up -d
	@until wget -O- http://localhost:4576/\?Action\=ListQueues >/dev/null 2>&1; do echo "Localstack is unreachable - sleeping"; sleep  1; done
	@echo "Creating SQS queues..."
	wget -qO- http://localhost:4576\?Action\=CreateQueue\&QueueName\=queue-test

dco-test-down:
	${DOCKERCOMPOSETEST} down