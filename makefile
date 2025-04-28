# Программа рассчитана на автоматический запуск на школьном маке или на wsl
# Просто набери make
# Программа автоматически установит Elastic в директорию /opt/goinfre/имя_пользователя/elasticsearch-номер_версии
# Если программа не запустится и будет ругаться на неподходящий пароль, просто сформируй новый пароль командой make change_password

OS:=$(shell uname -s)

LOCAL_PATH_ELASTIC:=/bin/elasticsearch
ELASTIC_PASS_RESET_ADDON:=-reset-password -u elastic -b --url "https://localhost:9200"

ifeq ($(OS), Darwin)
	PASSWORD_FILE:=elastic_pass_MAC.txt
	ELASTIC_DOWNLOAD_URL:=https://disk.yandex.ru/d/cy8KaNjG3u1p_w
	ELASTIC_DOWNLOAD_FILE:=elasticsearch-8.4.2-darwin-x86_64.tar.gz
	ELASTIC_NAME_DIR:=$(shell echo $(ELASTIC_DOWNLOAD_FILE) | sed -E 's/^([a-zA-Z0-9._]+-[0-9]+\.[0-9]+\.[0-9]+).*$$/\1/')
	ELASTIC_PATH:=/opt/goinfre/$(shell whoami)/$(ELASTIC_NAME_DIR)
else
	PASSWORD_FILE:=elastic_pass_WSL.txt
	ELASTIC_PATH:=$(shell cat elastic_dir_WSL.txt)
endif

ELASTIC_RUN:=$(ELASTIC_PATH)$(LOCAL_PATH_ELASTIC)
ELASTIC_INSTALL:=$(shell dirname $(ELASTIC_PATH))/$(ELASTIC_DOWNLOAD_FILE)
ELASTIC_LOGIN_PASS:=curl -XGET -k -u "elastic:$$(cat $(PASSWORD_FILE))"

all: elastic run_manual_parse

cleangocache:
	@go clean -cache -modcache -testcache && echo "\nОчистка кеша завершена успешно!\n"

elastic: $(ELASTIC_RUN)
	@curl -s -k "https://127.0.0.1:9200" >/dev/null 2>&1 && \
		echo "Elasticsearch уже запущен!" || \
		( \
			echo "Запускаем Elastic..."; \
			$(if $(filter $(OS),Darwin), \
				osascript -e 'tell application "Terminal" to do script "$(ELASTIC_RUN)"', \
				cmd.exe /c start wsl bash -c '$(ELASTIC_RUN); exec bash' \
			); \
			sleep 5; \
			while ! curl -s -k "https://127.0.0.1:9200" >/dev/null 2>&1; do sleep 2; done; \
			$(MAKE) --no-print-directory change_password; \
			echo "Elastic запущен!"; \
		)

$(ELASTIC_RUN):
	@echo "Файл Elastic не найден..."
	@$(MAKE) download
	@$(MAKE) extract

download:
	@echo "Скачиваем архив..."
	@curl -s 'https://cloud-api.yandex.net/v1/disk/public/resources/download?public_key=$(ELASTIC_DOWNLOAD_URL)' \
	| grep -o '"href":"[^"]*' \
	| grep -o 'https[^"]*' \
	| xargs curl -L -o $(ELASTIC_INSTALL)

extract:
	@echo "Распаковываем архив..."
	@tar -xzf $(ELASTIC_INSTALL) -C $(shell dirname $(ELASTIC_PATH))

change_password:
	@$(ELASTIC_RUN)$(ELASTIC_PASS_RESET_ADDON) | tail -n 1 | awk '{print $$NF}' | tr -d '\n' > $(PASSWORD_FILE)

test_elastic: elastic
	@$(ELASTIC_LOGIN_PASS) "https://localhost:9200/"

test_index: elastic
	@$(ELASTIC_LOGIN_PASS) "https://localhost:9200/places"

test_items: elastic
	@echo && echo ">>>> Объект  _id = 0: <<<<"
	@$(ELASTIC_LOGIN_PASS) "https://localhost:9200/places/_doc/0"
	@echo && echo && echo ">>>> Объект  _id = 13648: <<<<"
	@$(ELASTIC_LOGIN_PASS) "https://localhost:9200/places/_doc/13648" && echo

run_manual_parse: elastic
	@echo "Запуск ручного парсинга..."
	@sleep 2
	@cd src && go run -tags=manual .

run_auto_parse: elastic
	@echo "Запуск автоматического парсинга с помощью библиотеки gocsv..."
	@sleep 2
	@cd src && go run .


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