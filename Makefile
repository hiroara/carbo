.PHONY: compile
compile: internal/messaging/pb

.PHONY: clean
clean: clean-pb

.PHONY: recompile
recompile: clean compile

internal/messaging/pb:
	@mkdir -p internal/messaging/pb && docker compose run --rm protoc

.PHONY: clean-pb
clean-pb:
	@rm -rf internal/messaging/pb
