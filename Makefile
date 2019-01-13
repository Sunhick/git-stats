GO=go

TARGET=git-stats

all: ${TARGET}

build: ${TARGET}

run: ${TARGET}
	./${TARGET}

${TARGET}: src/git-stats.go
	${GO} build -o $@ $^

.PHONY: clean
clean:
	rm ${TARGET}
