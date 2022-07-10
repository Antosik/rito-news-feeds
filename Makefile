.PHONY: build

build:
	sam build \
		--template ./templates/sam.template.yaml \
		--parameter-overrides \
			ParameterKey=DomainName,ParameterValue=${DOMAIN_NAME} \
			ParameterKey=BucketName,ParameterValue=${BUCKET_NAME} \
			ParameterKey=DistributionId,ParameterValue=${DISTRIBUTION_ID}

local:
	sam local invoke "${FUNCTION}" \
		--parameter-overrides \
			ParameterKey=DomainName,ParameterValue=${DOMAIN_NAME} \
			ParameterKey=BucketName,ParameterValue=${BUCKET_NAME} \
			ParameterKey=DistributionId,ParameterValue=${DISTRIBUTION_ID}

validate:
	sam validate \
		--template ./templates/sam.template.yaml

deploy:
	sam deploy \
		--stack-name rito-news-feeds \
		--config-file ./templates/samconfig.toml \
		--parameter-overrides \
			ParameterKey=DomainName,ParameterValue=${DOMAIN_NAME} \
			ParameterKey=BucketName,ParameterValue=${BUCKET_NAME} \
			ParameterKey=DistributionId,ParameterValue=${DISTRIBUTION_ID}

cdn-create:
	aws cloudformation create-stack \
		--stack-name rito-news-cdn-stack \
		--template-body file://templates/cdn.template.yaml \
		--region us-east-1 \
		--capabilities CAPABILITY_NAMED_IAM \
		--parameters ParameterKey=DomainName,ParameterValue=${DOMAIN_NAME} ParameterKey=BucketName,ParameterValue=${BUCKET_NAME}

cdn-update:
	aws cloudformation update-stack \
		--stack-name rito-news-cdn-stack \
		--template-body file://templates/cdn.template.yaml \
		--region us-east-1 \
		--capabilities CAPABILITY_NAMED_IAM \
		--parameters ParameterKey=DomainName,ParameterValue=${DOMAIN_NAME} ParameterKey=BucketName,ParameterValue=${BUCKET_NAME}

build-LeagueOfLegendsStatusChecker:
	GOARCH=arm64 GOOS=linux go build -trimpath -o ./lolstatus lol/status/main.go lol/status/utils.go
	mv ./lolstatus $(ARTIFACTS_DIR)/bootstrap

build-LeagueOfLegendsNewsChecker:
	GOARCH=arm64 GOOS=linux go build -trimpath -o ./lolnews lol/news/main.go lol/news/utils.go
	mv ./lolnews $(ARTIFACTS_DIR)/bootstrap

build-LeagueOfLegendsEsportsChecker:
	GOARCH=amd64 GOOS=linux go build -trimpath -o ./lolesports lol/esports/main.go lol/esports/utils.go
	mv ./lolesports $(ARTIFACTS_DIR)/lolesports

build-VALORANTStatusChecker:
	GOARCH=arm64 GOOS=linux go build -trimpath -o ./valstatus val/status/main.go val/status/utils.go
	mv ./valstatus $(ARTIFACTS_DIR)/bootstrap

build-VALORANTNewsChecker:
	GOARCH=arm64 GOOS=linux go build -trimpath -o ./valnews val/news/main.go val/news/utils.go
	mv ./valnews $(ARTIFACTS_DIR)/bootstrap

build-VALORANTEsportsChecker:
	GOARCH=amd64 GOOS=linux go build -trimpath -o ./valesports val/esports/main.go val/esports/utils.go
	mv ./valesports $(ARTIFACTS_DIR)/valesports
