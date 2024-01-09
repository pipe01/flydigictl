PROTOS = $(wildcard *.proto)
PB_GOS = $(PROTOS:.proto=.pb.go)

.PHONY: all
all: $(PB_GOS)

%.pb.go: %.proto
	protoc --go_out=. --go_opt=paths=source_relative $^
