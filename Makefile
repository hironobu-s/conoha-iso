NAME=conoha-iso
BINDIR=bin
GOARCH=amd64

all: clean windows darwin linux

windows:
	GOOS=$@ GOARCH=$(GOARCH) CGO_ENABLED=0 go build $(GOFLAGS) -o $(BINDIR)/$@/$(NAME).exe
	cd bin/$@; zip $(NAME).$(GOARCH).zip $(NAME).exe

darwin:
	GOOS=$@ GOARCH=$(GOARCH) CGO_ENABLED=0 go build $(GOFLAGS) -o $(BINDIR)/$@/$(NAME)
	cd bin/$@; gzip -c $(NAME) > $(NAME)-osx.$(GOARCH).gz

linux:
	GOOS=$@ GOARCH=$(GOARCH) CGO_ENABLED=0 go build $(GOFLAGS) -o $(BINDIR)/$@/$(NAME)
	cd bin/$@; gzip -c $(NAME) > $(NAME)-linux.$(GOARCH).gz

clean:
	rm -rf $(BINDIR)

test:
	go test -v *.go
	go test -v command/*.go
