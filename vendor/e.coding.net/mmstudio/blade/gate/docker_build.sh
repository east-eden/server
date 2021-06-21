#!/bin/bash

BITBUCKET_TAG="$1"
IMAGE="zhongtai/gate"
LATEST_IMAGE="${IMAGE}:latest"
VERSION_IMAGE="${IMAGE}:${BITBUCKET_TAG}"
AWS_DEFAULT_REGION="us-east-1"
AWS_ECR_URI="public.ecr.aws/c4n2t7d7"
git checkout "${BITBUCKET_TAG}"
if [  $? -ne 0 ]
then
  echo "build error"
  exit 1
fi

pip3 install awscli
aws configure set aws_access_key_id "${AWS_ACCESS_KEY_ID}"
aws configure set aws_secret_access_key "${AWS_SECRET_ACCESS_KEY}"
aws --region "${AWS_DEFAULT_REGION}" ecr-public get-login-password | docker login --username AWS --password-stdin "${AWS_ECR_URI}"

export DOCKER_BUILDKIT=1 && docker build -t "${IMAGE}" --ssh=default="${HOME}/.ssh/id_rsa" -f Dockerfile.server .
docker tag "${LATEST_IMAGE}" "${AWS_ECR_URI}/${LATEST_IMAGE}"
docker tag "${LATEST_IMAGE}" "${AWS_ECR_URI}/${VERSION_IMAGE}"
#docker push "${AWS_ECR_URI}/${LATEST_IMAGE}"
#docker push "${AWS_ECR_URI}/${VERSION_IMAGE}"
aws configure set aws_access_key_id "${AWS_ACCESS_KEY_ID}"
aws configure set aws_secret_access_key "${AWS_SECRET_ACCESS_KEY}"
#aws --region "${AWS_DEFAULT_REGION}" ecr-public describe-images --repository-name "${IMAGE}"

