COVERDIR=$(CURDIR)/.cover
COVERAGEFILE=$(COVERDIR)/cover.out

test:
	@ginkgo --failFast ./...

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
	@until wget -O- http://localhost:9324/\?Action\=ListQueues >/dev/null 2>&1; do echo "ElasticMQ is unreachable - sleeping"; sleep  1; done
	@echo "Creating SQS queues..."
	wget -qO- http://localhost:9324\?Action\=CreateQueue\&QueueName\=queue-test

dcdn:
	@docker-compose down

.PHONY: test test-watch coverage coverage-ci coverage-html fmt vet dcup dcdn
