include ${STAGE}.env

#region SAM
build:
	sam build \
		--template ./templates/sam.template.yaml \
		--parameter-overrides \
			ParameterKey=DomainName,ParameterValue=${DOMAIN_NAME} \
			ParameterKey=BucketName,ParameterValue=${BUCKET_NAME} \
			ParameterKey=DistributionId,ParameterValue=${DISTRIBUTION_ID} \
			ParameterKey=Stage,ParameterValue=${STAGE}

local:
	sam local invoke "${FUNCTION}" \
		--parameter-overrides \
			ParameterKey=DomainName,ParameterValue=${DOMAIN_NAME} \
			ParameterKey=BucketName,ParameterValue=${BUCKET_NAME} \
			ParameterKey=DistributionId,ParameterValue=${DISTRIBUTION_ID} \
			ParameterKey=Stage,ParameterValue=${STAGE}

validate:
	sam validate \
		--template ./templates/sam.template.yaml \
		--lint

deploy-init:
	sam deploy \
		--stack-name rito-news-feeds-${STAGE} \
		--config-env ${STAGE} \
		--parameter-overrides \
			ParameterKey=DomainName,ParameterValue=${DOMAIN_NAME} \
			ParameterKey=BucketName,ParameterValue=${BUCKET_NAME} \
			ParameterKey=DistributionId,ParameterValue=${DISTRIBUTION_ID} \
			ParameterKey=Stage,ParameterValue=${STAGE} \
		--capabilities CAPABILITY_NAMED_IAM \
		--guided

deploy:
	sam deploy \
		--stack-name rito-news-feeds-${STAGE} \
		--config-file ./samconfig.toml \
		--config-env ${STAGE} \
		--parameter-overrides \
			ParameterKey=DomainName,ParameterValue=${DOMAIN_NAME} \
			ParameterKey=BucketName,ParameterValue=${BUCKET_NAME} \
			ParameterKey=DistributionId,ParameterValue=${DISTRIBUTION_ID} \
			ParameterKey=Stage,ParameterValue=${STAGE}

remove:
	sam delete \
		--stack-name rito-news-feeds-${STAGE} \
		--config-file ./templates/samconfig.toml \
		--config-env ${STAGE}
#endregion SAM

#region CDN
cdn-create:
	aws cloudformation create-stack \
		--stack-name rito-news-cdn-stack-${STAGE} \
		--template-body file://templates/cdn.template.yaml \
		--region us-east-1 \
		--capabilities CAPABILITY_NAMED_IAM \
		--parameters \
			ParameterKey=DomainName,ParameterValue=${DOMAIN_NAME} \
			ParameterKey=BucketName,ParameterValue=${BUCKET_NAME} \
			ParameterKey=Stage,ParameterValue=${STAGE}

cdn-update:
	aws cloudformation update-stack \
		--stack-name rito-news-cdn-stack-${STAGE} \
		--template-body file://templates/cdn.template.yaml \
		--region us-east-1 \
		--capabilities CAPABILITY_NAMED_IAM \
		--parameters \
			ParameterKey=DomainName,ParameterValue=${DOMAIN_NAME} \
		 	ParameterKey=BucketName,ParameterValue=${BUCKET_NAME} \
			ParameterKey=Stage,ParameterValue=${STAGE} 
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
	GOARCH=arm64 GOOS=linux go build -trimpath -o ./lornews lor/news/main.go lor/news/utils.go
	mv ./lornews $(ARTIFACTS_DIR)/bootstrap
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

build-WildRiftNewsChecker:
	GOARCH=arm64 GOOS=linux go build -trimpath -o ./wrnews wr/news/main.go wr/news/utils.go
	mv ./wrnews $(ARTIFACTS_DIR)/bootstrap
#endregion Build: WildRift

#endregion Build: RiotGames
build-RiotGamesNewsChecker:
	GOARCH=arm64 GOOS=linux go build -trimpath -o ./riotgamesnews riotgames/news/main.go riotgames/news/utils.go
	mv ./riotgamesnews $(ARTIFACTS_DIR)/bootstrap

build-RiotGamesJobsChecker:
	GOARCH=arm64 GOOS=linux go build -trimpath -o ./riotgamesjobs riotgames/jobs/main.go riotgames/jobs/utils.go
	mv ./riotgamesjobs $(ARTIFACTS_DIR)/bootstrap
#endregion Build: RiotGames
