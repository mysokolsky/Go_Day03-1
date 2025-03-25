OS := $(shell uname -s)

run_elastic:
	@echo "Starting Elasticsearch..."
	@if [ "$$(uname)" = "Darwin" ]; then \
		osascript -e 'tell application "Terminal" to do script  "$(shell cat elastic_run_MAC.txt)"'; \
	elif grep -qi microsoft /proc/version; then \
		wsl.exe -d Ubuntu -- bash -c "$(shell cat elastic_run_WSL.txt)"; \
	else \
		echo "Unsupported OS"; \
	fi

test_elastic:
ifeq ($(OS), Darwin)
	@bash elastic_test_MAC.txt
else
	@bash elastic_test_WSL.txt
endif

clean:
	@make removecash

include MAKEFILES/develop.mk
include MAKEFILES/git.mk
include MAKEFILES/clearmac.mk
include MAKEFILES/colors.mk