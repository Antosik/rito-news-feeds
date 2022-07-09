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
