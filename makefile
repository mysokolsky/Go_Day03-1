elastic_run:
	$(shell cat elastic_run_MAC.txt)
	$(shell cat elastic_run_WSL.txt)

elastic_test:
	$(shell cat elastic_test_MAC.txt)
	$(shell cat elastic_test_WSL.txt)

clean:
	@make removecash

include MAKEFILES/develop.mk
include MAKEFILES/git.mk
include MAKEFILES/clearmac.mk
include MAKEFILES/colors.mk