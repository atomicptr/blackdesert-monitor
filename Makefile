release: build
	zip -j blackdesert-monitor.zip build/*

build: clean
	mkdir -p build

	GOOS=windows go build -o build/blackdesert-monitor.exe
	cp settings.yaml.skel build/settings.yaml
	cp README.md build/README.md

clean:
	rm -rf build
