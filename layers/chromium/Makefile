HEADLESS_CHROMIUM_RELEASE=114.0.0

build-ChromiumLayer:
	# Create folders for binaries
	mkdir -p $(ARTIFACTS_DIR)/temp-chromium

	# Download and unzip headless chromium
	curl -SL https://github.com/Sparticuz/chromium/releases/download/v$(HEADLESS_CHROMIUM_RELEASE)/chromium-v$(HEADLESS_CHROMIUM_RELEASE)-pack.tar > $(ARTIFACTS_DIR)/headless-chromium.tar
	tar -xf $(ARTIFACTS_DIR)/headless-chromium.tar -C $(ARTIFACTS_DIR)/temp-chromium

	mkdir -p $(ARTIFACTS_DIR)/
	brotli -d $(ARTIFACTS_DIR)/temp-chromium/chromium.br -o '$(ARTIFACTS_DIR)/chromium'
	chmod 777 $(ARTIFACTS_DIR)/chromium

	rm $(ARTIFACTS_DIR)/headless-chromium.tar
	rm -rf $(ARTIFACTS_DIR)/temp-chromium
