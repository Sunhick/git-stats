# Copyright (c) 2019 Sunil

TARGET = git-stats

.PHONY: all
all: ${TARGET}

include common.mk

run: ${TARGET}
	./${TARGET}

${TARGET}: src/git-stats.go
	${GO} build -o $@ $^

.PHONY: clean
clean: decruft
	rm ${TARGET}
