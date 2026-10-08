FUNCTIONS := daily_schedule event_schedule event_notification invoice_reminder
CALL_TARGETS := $(addprefix call_,$(FUNCTIONS))

.DEFAULT_GOAL := help
.PHONY: build clean deploy $(CALL_TARGETS)

help:
	@echo "Available targets:"
	@printf "  %-23s - %s\n" "build" "Build Go binary for Linux"
	@printf "  %-23s - %s\n" "clean" "Remove built binary"
	@printf "  %-23s - %s\n" "deploy" "Clean, build, and deploy"
	@$(foreach fn,$(FUNCTIONS),printf "  %-23s - %s\n" "call_$(fn)" "Invoke deployed $(fn)";)
	@printf "  %-23s - %s\n" "help" "Show this help"
	@echo ""
	@echo "Note: call_* targets send a real email to the configured recipient."

build:
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bootstrap cmd/bootstrap/*.go

clean:
	rm -f ./bootstrap

deploy: clean build
	npx sls deploy --verbose

$(CALL_TARGETS): call_%:
	npx sls invoke --function $* --log
