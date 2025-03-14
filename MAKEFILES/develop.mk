# флаги для разработки
DEV_FLAGS=\
	-fsanitize=address \
	-pedantic
# -v

# переменная для разработки
CC_DEV=$(CXX) $(DEV_FLAGS)


# V_COMMAND=valgrind  --tool=memcheck --track-fds=yes --trace-children=yes --track-origins=yes --leak-check=full --show-leak-kinds=all -s
# L_COMMAND=leaks -atExit --


# цель для разработки - запуск просто make
compile: removecash
	clear
	@echo "///////////////////////////////////////////////////////////"
	@echo "/////////////////  компиляция  $(shell date +"%H:%M:%S")  //////////////////"
	@echo "///////////////////////////////////////////////////////////"
	@echo
	@rm -f run
	@$(CC_DEV) \
		DEVELOP/stack_test_manual.cc \
		-o \
		run
	@./run
# -cppcheck --enable=all --force	./*/*.cpp ./*.cpp ./*/*.hpp ./*.hpp


# ifeq ($(OS), Linux)
# 	valgrind --tool=memcheck --track-fds=yes --trace-children=yes --track-origins=yes --leak-check=full --show-leak-kinds=all -s  ./run
# else
# 	leaks -quiet -atExit -- ./run
# endif
# leaks -quiet -atExit -- ./run
# @rm -f run

# DEVELOP/*.c* \
# FUNCTIONS/*.c* \
# -lm \
# -lpcre \
# leaks -quiet -atExit -- ./run
# @rm -f run
# @rm -rf /../**.dSYM


# дополнительно
format: clean
	@cp ../materials/linters/.clang-format .
	-clang-format -style=Google -i ./*/*.cpp
	-clang-format -style=Google -i ./*.cpp
	-clang-format -style=Google -i ./*/*.hpp
	-clang-format -style=Google -i ./*.hpp

check: test
	@cp ../materials/linters/.clang-format .
	-clang-format -style=Google -n ./*/*.cpp
	-clang-format -style=Google -n ./*.cpp
	-clang-format -style=Google -n ./*/*.hpp
	-clang-format -style=Google -n ./*.hpp
# -sudo apt-get install -y cppcheck > /dev/null
# -brew install -y cppcheck > /dev/null
	-cppcheck --enable=all --force ./*/*.cpp ./*.cpp ./*/*.hpp ./*.hpp
	if [ `uname` = "Darwin" ]; \
	then \
		leaks --atExit -- ./test; \
	else \
		valgrind -s --leak-check=full --show-leak-kinds=all --track-origins=yes --tool=memcheck --log-file=valgrind.log ./test; \
	cat valgrind.log; \
	sleep 1; \
	fi
	@rm -f *.gcda *.gcno ./test



silent:
	@make all > /dev/null


# try: clean
# 	$(CC) -c -g DEVELOP/try_test.c
# 	ar rcs s21_matrix.a *.o
# 	ranlib s21_matrix.a
# 	sleep 1
# 	$(CC) DEVELOP/try_test.c s21_matrix.a -o try_test $(CHECK_FLAGS)
# 	valgrind -s --leak-check=full --show-leak-kinds=all --track-origins=yes --log-file=DEVELOP/valgrind_try_test.log ./try_test
# 	cat DEVELOP/valgrind_try_test.log
# 	rm -f try_test.o try_test s21_matrix.a
# 	sleep 1

