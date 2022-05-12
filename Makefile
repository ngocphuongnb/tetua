build_editor:
	echo "Building editor..."
	rm -rf packages/editor/dist/*
	cd packages/editor && yarn build

build_app:
	echo "Building app..."
	rm -rf views/*.go
	echo "package views" > views/views.go
	go run packages/prebuild/prebuild.go --force
	go run . bundlestatic
	go mod tidy

test_all:
	echo "Testing..."
	go test ./...

prerelease: build_app test_all

releaselocal: prerelease
	echo "Releasing..."
	rm -rf dist
	goreleaser release --snapshot --rm-dist

release: build_editor prerelease
	echo "Releasing..."
	rm -rf dist
	goreleaser release