build:
	mkdir -p build/
	go build -o build/chirpstack-gw-protobuf-translator app.go

clean:
	rm -rf build