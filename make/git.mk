GIT_COMMIT_ID := $(shell git rev-parse HEAD)
ifneq ($(shell git status --porcelain --untracked-files=no),)
       GIT_COMMIT_ID := latest
endif

GIT_COMMIT_ID_SHORT := $(shell git rev-parse --short HEAD)
ifneq ($(shell git status --porcelain --untracked-files=no),)
       GIT_COMMIT_ID_SHORT := latest
endif

BUILD_TIME = `date -u '+%Y-%m-%dT%H:%M:%SZ'`