
killall -9 McmAPI
rm -f ./McmAPI
go build McmAPI.go DBProcess.go FilterProcess.go
./McmAPI &disown 

#go build McmAPI.go
#time ./McmAPI
