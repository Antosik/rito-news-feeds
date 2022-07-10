HEADLESS_CHROMIUM_RELEASE = 1.0.0-57

build-ChromiumLayer:
	# Create folders for binaries
	mkdir -p $(ARTIFACTS_DIR)/bin

	# Download and unzip headless chromium
	curl -SL https://github.com/adieuadieu/serverless-chrome/releases/download/v$(HEADLESS_CHROMIUM_RELEASE)/stable-headless-chromium-amazonlinux-2.zip > $(ARTIFACTS_DIR)/headless-chromium.zip
	unzip $(ARTIFACTS_DIR)/headless-chromium.zip -d $(ARTIFACTS_DIR)/bin/
	mv $(ARTIFACTS_DIR)/bin/headless-chromium $(ARTIFACTS_DIR)/bin/chromium
	rm $(ARTIFACTS_DIR)/headless-chromium.zip
