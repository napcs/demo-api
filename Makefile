version = 0.1.0

mac:
	mkdir -p dist/mac
	cp data.json.example dist/mac/data.json
	cp README.md dist/mac/README
	cp LICENSE dist/mac/LICENSE
	env GOOS=darwin GOARCH=amd64 go build -o dist/mac/demo-api
	cd dist/mac; zip -9 demo-api-macos-${version}.zip data.json demo-api README LICENSE

linux:
	mkdir -p dist/linux
	cp README.md dist/linux/README
	cp LICENSE dist/linux/LICENSE
	cp data.json.example dist/linux/data.json
	env GOOS=linux GOARCH=amd64 go build -o dist/linux/demo-api
	cd dist/linux; zip -9 demo-api-linux-${version}.zip data.json demo-api README LICENSE

windows:
	mkdir -p dist/windows
	cp README.md dist/windows/README
	cp LICENSE dist/windows/LICENSE
	cp data.json.example dist/windows/data.json
	env GOOS=windows GOARCH=amd64 go build -o dist/windows/demo-api.exe
	cd dist/windows; zip -9 demo-api-windows-${version}.zip data.json demo-api.exe README LICENSE

clean:
	rm -rf dist/
