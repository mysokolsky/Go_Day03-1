
GITBRANCH:=$(shell git rev-parse --abbrev-ref HEAD)

GITHUBPROJECTNAME:=$(shell basename $(shell pwd))

GITHUB_TOKEN:=$(shell cat github_token.tok)

GITLABPROJECTURL:=$(shell git remote get-url origin)

add:
	git add .; sleep 1

commit:
	@git commit -m "======= $(GITBRANCH) === $(shell date +'ДАТА %d-%m-%y === ВРЕМЯ %H:%M:%S') ======="; sleep 1

# автоматическая заливка на гит-репозиторий с автоматическим коммитом по текущей дате и времени
push: add commit
	git push origin $(GITBRANCH)

# автоматическая загрузка с гит-репозитория на текущую машину
pull:
	git stash
	git pull origin $(GITBRANCH)

init:
	git init
	git add .
	git commit -m "init"

new-branch-develop:
	git checkout -b develop

create-repo:
	curl -u 'mysokolsky:$(GITHUB_TOKEN)' https://api.github.com/user/repos -d'{"name":"$(GITHUBPROJECTNAME)"}'

# сначала уверждение гитлаб-репозитория как основного синхронизируемого, а потом подключение дополнительного для закачки репозитория на гитхаб
add-remote-repo:
	git remote add origin $(GITLABPROJECTURL)
	git remote set-url --add --push origin git@github.com:mysokolsky/$(GITHUBPROJECTNAME).git
	git remote -v
# обнуление пуш-адресов
# git remote set-url --push origin ""
# проверка
# git remote get-url --push origin