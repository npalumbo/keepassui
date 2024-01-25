#@ Helpers
# from https://www.thapaliya.com/en/writings/well-documented-makefiles/
help:  ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Tools
tools: ## Installs required binaries locally
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install fyne.io/fyne/v2/cmd/fyne@latest
	go install go.uber.org/mock/mockgen@latest

create-android-tools-dir:
	mkdir -p ~/tools

install-android-ndk-25:
	wget -O ~/tools/android-ndk-r25c.zip https://dl.google.com/android/repository/android-ndk-r25c-linux.zip  && unzip ~/tools/android-ndk-r25c.zip -d ~/tools

install-android-cmdline-tools:
	mkdir -p ~/tools/android-sdk && wget -O ~/tools/android-cmdline-tools.zip https://dl.google.com/android/repository/commandlinetools-linux-10406996_latest.zip && unzip ~/tools/android-cmdline-tools.zip -d ~/tools/android-sdk/ && mv ~/tools/android-sdk/cmdline-tools ~/tools/android-sdk/latest && mkdir ~/tools/android-sdk/cmdline-tools && mv ~/tools/android-sdk/latest ~/tools/android-sdk/cmdline-tools

install-android-platform-tools:
	cd ~/tools/android-sdk && ANDROID_SDK_ROOT=~/tools/android-sdk ~/tools/android-sdk/cmdline-tools/latest/bin/sdkmanager --sdk_root=$$ANDROID_SDK_ROOT --install "platform-tools"

install-android-build-tools:
	cd ~/tools/android-sdk/cmdline-tools/latest/bin && ANDROID_SDK_ROOT=~/tools/android-sdk ./sdkmanager --sdk_root=$$ANDROID_SDK_ROOT --install  "build-tools;34.0.0" && mv build-tools ~/tools/android-sdk/

install-android-bundletool:
	mkdir -p ~/bin && mkdir -p ~/tools/bundletools && wget -O ~/tools/bundletools/bundletool-all.jar https://github.com/google/bundletool/releases/download/1.15.6/bundletool-all-1.15.6.jar
	@echo '#!/usr/bin/bash' > ~/bin/bundletool
	@echo 'java -jar ~/tools/bundletools/bundletool-all.jar "$$@"' >> ~/bin/bundletool
	@chmod +x ~/bin/bundletool


##@ Install All Android tools
install-all-android-tools: create-android-tools-dir install-android-ndk-25 install-android-cmdline-tools install-android-platform-tools install-android-build-tools install-android-bundletool ## Install all the required Android tools


##@ Create Android Package (APK file)
package-android: ## Generates an android APK file
	ANDROID_NDK_HOME=~/tools/android-ndk-r25c fyne package -os android

##@ Create Android Release (AAB file)
release-android: ## Generates an android AAB file
	ANDROID_NDK_HOME=~/tools/android-ndk-r25c ANDROID_HOME=~/tools/android-sdk fyne release -os android -keyStore ~/dev/gplay.keystore -keyName alias

##@ Building
build-multi-arch: ## Builds keepassui go binary for linux and darwin. Outputs to `bin/keepassui-$GOOS-$GOARCH`.
	@echo "== build-multi-arch"
	mkdir -p bin/
	GOOS=linux GOARCH=amd64 CGO_ENABL	ED=0 go build -o bin/keepassui-linux-amd64 ./...
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o bin/keepassui-darwin-amd64 ./...
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o bin/keepassui-darwin-arm64 ./...

build: check## Builds keepassui go binary for local arch. Outputs to `bin/keepassui`
	@echo "== build"
	CGO_ENABLED=1 go build -o bin/ ./...

##@ Cleanup
clean: ## Deletes binaries from the bin folder
	@echo "== clean"
	rm -rfv ./bin

##@ Tests
test: ## Run unit tests
	@echo "== unit test"
	go test ./...

##@ Run static checks
check: ## Runs lint, fmt and vet checks against the codebase
	golangci-lint --timeout 120s run
	go fmt ./...
	go vet ./...

##@ Golang Generate
generate: ## Calls golang generate
	go mod tidy
	go generate ./...

