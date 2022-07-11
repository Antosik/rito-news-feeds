.PHONY: build

#region SAM
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
#endregion SAM

#region CDN
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
#endregion CDN

#region Build: League of Legends
build-LeagueOfLegendsStatusChecker:
	GOARCH=arm64 GOOS=linux go build -trimpath -o ./lolstatus lol/status/main.go lol/status/utils.go
	mv ./lolstatus $(ARTIFACTS_DIR)/bootstrap

build-LeagueOfLegendsNewsChecker:
	GOARCH=arm64 GOOS=linux go build -trimpath -o ./lolnews lol/news/main.go lol/news/utils.go
	mv ./lolnews $(ARTIFACTS_DIR)/bootstrap

build-LeagueOfLegendsEsportsChecker:
	GOARCH=amd64 GOOS=linux go build -trimpath -o ./lolesports lol/esports/main.go lol/esports/utils.go
	mv ./lolesports $(ARTIFACTS_DIR)/lolesports
#endregion Build: League of Legends

#region Build: VALORANT
build-VALORANTStatusChecker:
	GOARCH=arm64 GOOS=linux go build -trimpath -o ./valstatus val/status/main.go val/status/utils.go
	mv ./valstatus $(ARTIFACTS_DIR)/bootstrap

build-VALORANTNewsChecker:
	GOARCH=arm64 GOOS=linux go build -trimpath -o ./valnews val/news/main.go val/news/utils.go
	mv ./valnews $(ARTIFACTS_DIR)/bootstrap

build-VALORANTEsportsChecker:
	GOARCH=amd64 GOOS=linux go build -trimpath -o ./valesports val/esports/main.go val/esports/utils.go
	mv ./valesports $(ARTIFACTS_DIR)/valesports
#endregion Build: VALORANT

#region Build: Legends of Runeterra
build-LegendsOfRuneterraStatusChecker:
	GOARCH=arm64 GOOS=linux go build -trimpath -o ./lorstatus lor/status/main.go lor/status/utils.go
	mv ./lorstatus $(ARTIFACTS_DIR)/bootstrap

build-LegendsOfRuneterraNewsChecker:
	GOARCH=amd64 GOOS=linux go build -trimpath -o ./lornews lor/news/main.go lor/news/utils.go
	mv ./lornews $(ARTIFACTS_DIR)/lornews
#endregion Build: Legends of Runeterra

#region Build: Teamfight Tactics
build-TeamfightTacticsNewsChecker:
	GOARCH=arm64 GOOS=linux go build -trimpath -o ./tftnews tft/news/main.go tft/news/utils.go
	mv ./tftnews $(ARTIFACTS_DIR)/bootstrap
#endregion Build: Teamfight Tactics

#region Build: WildRift
build-WildRiftStatusChecker:
	GOARCH=arm64 GOOS=linux go build -trimpath -o ./wrstatus wr/status/main.go wr/status/utils.go
	mv ./wrstatus $(ARTIFACTS_DIR)/bootstrap
#endregion Build: WildRift
