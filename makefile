OS := $(shell uname -s)

LOCAL_PATH_ELASTIC:=/bin/elasticsearch
ELASTIC_PASS_RESET_ADDON:=-reset-password -u elastic -b

ifeq ($(OS), Darwin)
	PASSWORD_FILE:=elastic_pass_MAC.txt
	ELASTIC_PATH:=$(shell cat elastic_dir_MAC.txt)
else
	PASSWORD_FILE:=elastic_pass_WSL.txt
	ELASTIC_PATH:=$(shell cat elastic_dir_WSL.txt)
endif

run_elastic: change_password
ifeq ($(OS), Darwin)
	@osascript -e 'tell application "Terminal" to do script  "$(ELASTIC_PATH)$(LOCAL_PATH_ELASTIC)"'
else
	@cmd.exe /c start wsl bash -c '$(ELASTIC_PATH)$(LOCAL_PATH_ELASTIC); exec bash'
endif

change_password:
	$(ELASTIC_PATH)$(LOCAL_PATH_ELASTIC)$(ELASTIC_PASS_RESET_ADDON) | tail -n 1 | awk '{print $$NF}' | tr -d '\n' > $(PASSWORD_FILE)


test_elastic: change_password
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