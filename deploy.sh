#!/bin/bash

BUILD_JOB_ID=$(cat build_job_id)
echo "Build Job ID: ${BUILD_JOB_ID}"
ARTIFACT_URL="${CI_PROJECT_URL}/-/jobs/${BUILD_JOB_ID}/artifacts/download?archive=zip"
echo "Artifact URL: ${ARTIFACT_URL}"
levant deploy \
  -var git_sha="${CI_COMMIT_SHORT_SHA}" \
  -var artifact_url="${ARTIFACT_URL}" \
  webapp-prod.nomad
