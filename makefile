OS := $(shell uname -s)

run_elastic:
ifeq ($(OS), Darwin)
	@osascript -e 'tell application "Terminal" to do script  "$(shell cat elastic_run_MAC.txt)"'
else
	@cmd.exe /c start wsl bash -c '$(shell cat elastic_run_WSL.txt); exec bash'
endif

test_elastic:
ifeq ($(OS), Darwin)
	@bash elastic_test_MAC.txt
else
	@bash elastic_test_WSL.txt
endif

gomodinit:
	go mod init $(shell basename $(PWD))

gogetelastic:
	go get github.com/elastic/go-elasticsearch/v8

clean:
	@make removecash

include MAKEFILES/develop.mk
include MAKEFILES/git.mk
include MAKEFILES/clearmac.mk
include MAKEFILES/colors.mk