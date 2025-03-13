newrepo:
	curl -u 'mysokolsky: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' https://api.github.com/user/repos -d'{"name":"GO/Go_Day03"}'
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