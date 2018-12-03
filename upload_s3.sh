#!/bin/bash
BUCKET=`sceptre --output json describe-stack-outputs example bucket | jq -r '.[] | select(.OutputKey=="BucketName") | .OutputValue'`
aws s3 cp dist/lambda.zip s3://$BUCKET/lambda.zip
