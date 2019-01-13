# Copyright (c) 2019 Sunil

TARGET = git-stats
SRCS = src/git-stats.go

.PHONY: all
all: ${TARGET}

include common.mk

rebuild: clean ${TARGET}

run: ${TARGET}
	./${TARGET}

${TARGET}: ${SRCS}
	${E} "go - " $@
	${Q} ${GO} build -o $@ $^

.PHONY: clean
clean: decruft
	rm ${TARGET}
