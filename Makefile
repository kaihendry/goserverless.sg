STACK = gosls
VERSION = $(shell git rev-parse --abbrev-ref HEAD)-$(shell git rev-parse --short HEAD)

.PHONY: build deploy validate destroy

DOMAINNAME = gosls.dabase.com
ACMCERTIFICATEARN = arn:aws:acm:ap-southeast-1:407461997746:certificate/87b0fd84-fb44-4782-b7eb-d9c7f8714908

deploy:
	sam build
	sam deploy --resolve-s3 --stack-name $(STACK) --parameter-overrides DomainName=$(DOMAINNAME) ACMCertificateArn=$(ACMCERTIFICATEARN) --no-confirm-changeset --no-fail-on-empty-changeset --capabilities CAPABILITY_IAM

build-MainFunction:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=${VERSION}" -o ${ARTIFACTS_DIR}/main main.go

validate:
	aws cloudformation validate-template --template-body file://template.yml

destroy:
	aws cloudformation delete-stack --stack-name $(STACK)

sam-tail-logs:
	sam logs --stack-name $(STACK) --tail
