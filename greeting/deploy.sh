#!/bin/bash

PROJECT_ID=instruments-cloud-function

function deploy() {
  local function_name=$1
  gcloud functions deploy ${function_name} \
  --gen2 \
  --runtime=go121 \
  --region=asia-northeast1 \
  --project=${PROJECT_ID} \
  --source=. \
  --entry-point=Greeting \
  --trigger-http \
  --allow-unauthenticated \
  --env-vars-file=env-${function_name}.yaml
}

deploy greeting-1
deploy greeting-2