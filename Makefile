all: compliance gofumpt staticcheck

compliance:
	# not using `go run` here because go-licenses doesn't seem to want to run with go 1.22
    #
    # ignoring 'github.com/jmespath/go-jmespath' because of unrecognizable license (a mention
    # of Apache-2.0 at the 0.4.0 commit
	go-licenses save   --ignore github.com/jmespath/go-jmespath --ignore github.com/chrismarget --save_path Third_Party_Code --force ./... || exit 1 ;\
	go-licenses report --ignore github.com/jmespath/go-jmespath --ignore github.com/chrismarget --template .notices.tpl ./... > Third_Party_Code/NOTICES.md || exit 1 ;\

	## workaround for github.com/jmespath/go-jmespath license
	mkdir -p Third_Party_Code/github.com/jmespath/go-jmespath ;\
    curl -s 'https://raw.githubusercontent.com/jmespath/go-jmespath/v0.4.0/LICENSE' >> Third_Party_Code/github.com/jmespath/go-jmespath/LICENSE ;\
	sh -c '( \
      echo "" ;\
      echo "## github.com/jmespath/go-jmespath" ;\
      echo "" ;\
      echo "* Name: github.com/jmespath/go-jmespath" ;\
      echo "* Version: v0.4.0" ;\
      echo "* License: [Apache-2.0](https://github.com/jmespath/go-jmespath/blob/v0.4.0/LICENSE)" ;\
      echo "" ;\
	  echo "\`\`\`" ;\
	  cat Third_Party_Code/github.com/jmespath/go-jmespath/LICENSE ;\
	  echo "\`\`\`" ;\
      echo "" ;\
	)' >> Third_Party_Code/NOTICES.md ;\

gofumpt:
	go run mvdan.cc/gofumpt -w src

staticcheck:
	go run honnef.co/go/tools/cmd/staticcheck ./...
