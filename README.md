# GCSC (Google Cloud Snapshots Cleaner)

GCSC allows you to clean up disk snapshots in Google Cloud based on some rules you set.

[![Build Status](https://travis-ci.org/Fale/gcsc.svg?branch=master)](https://travis-ci.org/Fale/gcsc)

[![Go Report Card](https://goreportcard.com/badge/github.com/fale/gcsc)](https://goreportcard.com/report/github.com/fale/gcsc)

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

## Configuration file
You should create a configuration file in `~/.gcsc/config.yaml` with a content similar to this:

    project-id: your-gcp-project-id
    retention-policies:
      - begin: 0h
        end: 168h # 7 days
        cadence: 1h
      - begin: 168h # 7 days
        end: 336h # 14 days
        cadence: 24h
      - begin: 336h # 14 days
        end: 1512h # 63 days
        cadence: 168h # 7 days
      - begin: 1512h # 63 days
        end: 25512h # 1063 days
        cadence: 24000h # 1000 days

## Running the application
You can run the application with:

    ./gcsc

## Usage

    Usage:
      snapshot-cleaner [command]
    
    Available Commands:
      clean       execute a cleaning
      help        Help about any command
      http        listen to HTTP port

    Flags:
          --automatic           Include automatic backups (default true)
          --dry-run             Dry run mode
      -h, --help                help for snapshot-cleaner
          --manual              Include manual backups
      -p, --project-id string   Google Cloud Project ID

    Use "snapshot-cleaner [command] --help" for more information about a command.
