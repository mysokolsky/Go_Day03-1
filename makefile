OS := $(shell uname -s)

elastic_run:
ifeq ($(OS), Darwin)
	@bash elastic_run_MAC.txt
else
	@bash elastic_run_WSL.txt
endif

elastic_test:
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