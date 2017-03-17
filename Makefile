GOBIN = $(value GOPATH)/bin
.PHONY: test deps gx work pub
# Anything which requires deps should end with: gx-go rewrite --undo

hc: deps
	go install ./cmd/hc
	gx-go rewrite --undo
bs: deps
	go install ./cmd/bs
	gx-go rewrite --undo
init: hc
	hc init node@example.com
test: testcore
	gx-go rewrite --undo
testcore: deps
	go get -t
	go test -v ./...||exit 1
testall: testcore hc init
	hc --debug --verbose clone --force examples/chat   examples-chat   && hc --debug --verbose test examples-chat
	hc --debug --verbose clone --force examples/sample examples-sample && hc --debug --verbose test examples-sample
deps: gx
	gx-go rewrite
	go get -d ./...
gx: $(GOBIN)/gx $(GOBIN)/gx-go
	gx install --global
$(GOBIN)/gx:
	go get -u github.com/whyrusleeping/gx
$(GOBIN)/gx-go:
	go get -u github.com/whyrusleeping/gx-go
work: $(GOBIN)/gx-go
	gx-go rewrite
pub: $(GOBIN)/gx-go
	gx-go rewrite --undo
