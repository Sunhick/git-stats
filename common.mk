# Copyright (c) 2019 Sunil

GO = go
Q = @
E = @echo

CRUFT = $(wildcard *~)

.PHONY: decruft
decruft:
	${Q} ${RM} -- ${CRUFT}
