# Google Cloud Snapshots Cleaner

Google Cloud Snapshots Cleaner allows you to clean up disk snapshots in Google Cloud based on some rules you set.

[![Build Status](https://travis-ci.org/Fale/google-cloud-snapshot-cleaner.svg?branch=master)](https://travis-ci.org/kelfa/kelfa)

[![Go Report Card](https://goreportcard.com/badge/github.com/fale/google-cloud-snapshot-cleaner)](https://goreportcard.com/report/github.com/fale/google-cloud-snapshot-cleaner)

## Create credentials
In some cases credentials are not needed (ie: if the program runs on an GCE instance with appropriate Service Account).
In some cases credentials are needed.
When they are, they can be created with:

    PROJECT_ID=$(gcloud config list --format 'value(core.project)')
    gcloud iam roles create snapshot-cleaner --project ${PROJECT_ID} --file role.yaml
    gcloud beta iam service-accounts create snapshot_cleaner --display-name "Snapshot Cleaner"
    gcloud projects add-iam-policy-binding ${PROJECT_ID} --member serviceAccount:snapshot-cleaner@${PROJECT_ID}.iam.gserviceaccount.com --role projects/ed-aem6/roles/snapshot_cleaner
    gcloud iam service-accounts keys create ~/key.json --iam-account snapshot-cleaner@${PROJECT_ID}.iam.gserviceaccount.com

If you are using a JSON credential, remember to export the `GOOGLE_APPLICATION_CREDENTIALS` variable:

    export GOOGLE_APPLICATION_CREDENTIALS=~/key.json

## Running the application
You can run the application with:

    ./google-cloud-snapshot-cleaner

## Usage

    Usage:
      snapshot-cleaner [flags]
    
    Flags:
          --automatic           Include automatic backups (default true)
          --dry-run             Dry run mode
      -h, --help                help for snapshot-cleaner
          --manual              Include manual backups
      -p, --project-id string   Google Cloud Project ID
