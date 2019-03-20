normal: clean
	go build

# Generate code for linux
linux: clean
	env GOOS=linux GOARCH=amd64 go build

# Generate code for windows
windows:
	env GOOS=windows GOARCH=amd64 go build

# Clean all old files
clean:
	rm -f calculator
	rm -f *.exe
