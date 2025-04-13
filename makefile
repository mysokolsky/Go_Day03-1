OS:=$(shell uname -s)

LOCAL_PATH_ELASTIC:=/bin/elasticsearch
ELASTIC_PASS_RESET_ADDON:=-reset-password -u elastic -b --url "https://localhost:9200"

ifeq ($(OS), Darwin)
	PASSWORD_FILE:=elastic_pass_MAC.txt
	ELASTIC_PATH:=$(shell cat elastic_dir_MAC.txt)
	ELASTIC_DOWNLOAD_URL:=https://downloader.disk.yandex.ru/disk/b22dec595f91c451745f7e1f8dca95d5029353de6fa53f9ec4ebdc2ca51bc7df/67f6f169/fKqInKw3d7bLFOeFnMGnhNoBog5PxpbChttBBDC3x3UsD4eAa5YILQcvwb2Ms_j8vtAD8IWpYqSgCCXV75S9g9-mtHtQbq1XFrfoMIHCY3Sr8npumZHI4midPdWhecNq?uid=0&filename=elasticsearch-8.4.2-darwin-x86_64.tar.gz&disposition=attachment&hash=xhdH749L3BMuzoZV8bO6PHdqqp8rWmRmo0pFHuv781sNePa3sfRUD8Sh%2BPb7xhJeskEAmkQ4kXOg0TR8ZsXayQ%3D%3D&limit=0&content_type=application%2Fx-gzip&owner_uid=1130000030982078&fsize=379332642&hid=ecb9826972a8625a38039a6a50f6e649&media_type=compressed&tknv=v2
	ELASTIC_DOWNLOAD_FILE:=elasticsearch-8.4.2-darwin-x86_64.tar.gz
else
	PASSWORD_FILE:=elastic_pass_WSL.txt
	ELASTIC_PATH:=$(shell cat elastic_dir_WSL.txt)
endif

ELASTIC_RUN:=$(ELASTIC_PATH)$(LOCAL_PATH_ELASTIC)
ELASTIC_INSTALL:=$(shell dirname $(ELASTIC_PATH))/$(ELASTIC_DOWNLOAD_FILE)

all: run

run: $(ELASTIC_RUN)
	@echo "Запускаем..."
ifeq ($(OS), Darwin)
	@osascript -e 'tell application "Terminal" to do script  "$(ELASTIC_RUN)"'
else
	@cmd.exe /c start wsl bash -c '$(ELASTIC_RUN); exec bash'
endif
	@sleep 5 && while ! curl -s -k "https://127.0.0.1:9200" >/dev/null; do sleep 2; done
	@$(MAKE) --no-print-directory change_password

$(ELASTIC_RUN):
	@echo "Файл не найден..."
	$(MAKE) download
	$(MAKE) extract

download:
	@echo "Скачиваем архив..."
	curl -L -o $(ELASTIC_INSTALL) $(ELASTIC_DOWNLOAD_URL)

extract:
	@echo "Распаковываем архив..."
	tar -xzf $(ELASTIC_INSTALL) -C $(shell dirname $(ELASTIC_PATH))

change_password:
	@$(ELASTIC_PATH)$(LOCAL_PATH_ELASTIC)$(ELASTIC_PASS_RESET_ADDON) | tail -n 1 | awk '{print $$NF}' | tr -d '\n' > $(PASSWORD_FILE)

test_elastic:
	@curl -XGET -k -u "elastic:$$(cat $(PASSWORD_FILE))" "https://localhost:9200/"

test_index:
	@curl -XGET -k -u "elastic:$$(cat $(PASSWORD_FILE))" "https://localhost:9200/places"

test_items:
	@echo && echo ">>>> Объект  _id = 0: <<<<"
	@curl -XGET -k -u "elastic:$$(cat $(PASSWORD_FILE))" "https://localhost:9200/places/_doc/0"
	@echo && echo && echo ">>>> Объект  _id = 13648: <<<<"
	@curl -XGET -k -u "elastic:$$(cat $(PASSWORD_FILE))" "https://localhost:9200/places/_doc/13648" && echo

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