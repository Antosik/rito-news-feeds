.PHONY: build

build:
	sam build \
		--template ./templates/sam.template.yaml \
		--parameter-overrides BucketName=${BUCKET_NAME}

local:
	sam local invoke "${FUNCTION}"
		--parameter-overrides BucketName=${BUCKET_NAME}

cdn-create:
	aws cloudformation create-stack \
		--stack-name rito-news-cdn-stack \
		--template-body file://templates/cdn.template.yaml \
		--region us-east-1 \
		--capabilities CAPABILITY_NAMED_IAM
		--parameters DomainName=${DOMAIN_NAME},BucketName=${BUCKET_NAME}

cdn-update:
	aws cloudformation update-stack \
		--stack-name rito-news-cdn-stack \
		--template-body file://templates/cdn.template.yaml \
		--region us-east-1 \
		--capabilities CAPABILITY_NAMED_IAM
		--parameters DomainName=${DOMAIN_NAME},BucketName=${BUCKET_NAME}
