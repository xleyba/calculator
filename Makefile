normal: clean
	go build

# Generate code for linux
linux: clean
	env GOOS=linux GOARCH=amd64 go build
	zip -9 calculator_linux_v1.0.zip calculator conf.yaml

# Generate code for windows
windows: clean
	env GOOS=windows GOARCH=amd64 go build
	mv calculator.exe calculator_amd64.exe
	zip -9 calculator_win_v1.0.zip calculator_amd64.exe conf.yaml

# Clean all old files
clean:
	rm -f calculator
	rm -f *.exe
