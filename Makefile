.PHONY: compile
compile: pb

.PHONY: clean
clean: clean-pb

.PHONY: recompile
recompile: clean compile

pb:
	@mkdir -p pb && docker compose run --rm protoc

.PHONY: clean-pb
clean-pb:
	@rm -rf pb
