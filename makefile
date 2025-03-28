OS := $(shell uname -s)

ifeq ($(OS), Darwin)
	PASSWORD_FILE:=elastic_pass_MAC.txt
else
	PASSWORD_FILE:=elastic_pass_WSL.txt
endif

run_elastic:
ifeq ($(OS), Darwin)
	@osascript -e 'tell application "Terminal" to do script  "$(shell cat elastic_run_MAC.txt)"'
else
	@cmd.exe /c start wsl bash -c '$(shell cat elastic_run_WSL.txt); exec bash'
endif

test_elastic:
	curl -XGET -k -u "elastic:$$(cat $(PASSWORD_FILE))" "https://localhost:9200/"

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