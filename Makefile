
default:
	go mod tidy -v
	go test -trimpath=false -race=false -v
	-go test -v -run=. -trimpath=false -race=false -coverprofile tracer.out ./...
	-go tool cover -html=tracer.out -o tracer.html
	if [ -x ./cover-publish ] ; then ./cover-publish tracer.html ; fi

clean:
	go clean -v
	rm -f *~ tracer.out tracer.html

