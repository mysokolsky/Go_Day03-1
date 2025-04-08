OS:=$(shell uname -s)

LOCAL_PATH_ELASTIC:=/bin/elasticsearch
ELASTIC_PASS_RESET_ADDON:=-reset-password -u elastic -b --url "https://localhost:9200"

ifeq ($(OS), Darwin)
	PASSWORD_FILE:=elastic_pass_MAC.txt
	ELASTIC_PATH:=$(shell cat elastic_dir_MAC.txt)
else
	PASSWORD_FILE:=elastic_pass_WSL.txt
	ELASTIC_PATH:=$(shell cat elastic_dir_WSL.txt)
endif

run_elastic:
ifeq ($(OS), Darwin)
	@osascript -e 'tell application "Terminal" to do script  "$(ELASTIC_PATH)$(LOCAL_PATH_ELASTIC)"'
else
	@cmd.exe /c start wsl bash -c '$(ELASTIC_PATH)$(LOCAL_PATH_ELASTIC); exec bash'
endif
	@sleep 5 && while ! curl -s -k "https://127.0.0.1:9200" >/dev/null; do sleep 2; done
	@$(MAKE) --no-print-directory change_password

change_password:
	@$(ELASTIC_PATH)$(LOCAL_PATH_ELASTIC)$(ELASTIC_PASS_RESET_ADDON) | tail -n 1 | awk '{print $$NF}' | tr -d '\n' > $(PASSWORD_FILE)

test_elastic:
	@curl -XGET -k -u "elastic:$$(cat $(PASSWORD_FILE))" "https://localhost:9200/"

test_index:
	curl -XGET -k -u "elastic:$$(cat $(PASSWORD_FILE))" "https://localhost:9200/places"

test_item:
	curl -XGET -k -u "elastic:$$(cat $(PASSWORD_FILE))" "https://localhost:9200/places/_search?scroll=1m" -H 'Content-Type: application/json' -d'{"size": 3,"query": {"match_all": {}}}'

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