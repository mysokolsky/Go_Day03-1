newrepo:
	curl -u 'mysokolsky:TOKEN' https://api.github.com/user/repos -d'{"name":"GO/Go_Day03"}'
	git init
	git checkout -b develop
	git add .
	git commit -m "first"
	git remote add origin git@github.com:mysokolsky/GO/Go_Day03.git
	git push origin develop

push:
	git add .
	git commit -m "now"
	git push origin develop