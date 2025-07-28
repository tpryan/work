clean: 
	-rm collect
	-rm go.sum


collect: clean go.sum
	go build -o collect . exec
	chmod +x collect 

go.sum: clean
	go mod tidy	