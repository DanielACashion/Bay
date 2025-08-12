# Bay
Meant to be the go code that serves as a proof of concept for language agnostic program hotreloader with protobufs for defined contracts between programming langs


# Build
made using go 1.20.4, root folder 
go build .
go build -o Bay.dll -buildmode=c-shared main.go

# Todo
currently determining whether this will use jni or some popular package instead
(currently evaluating pros and cods of jni vs jna)
