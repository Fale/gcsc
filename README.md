mkdir -p ~/.snapshot-cleaner

PROJECT_ID=$(gcloud config list --format 'value(core.project)')
gcloud iam roles create snapshot-cleaner --project ${PROJECT_ID} --file role.yaml
gcloud beta iam service-accounts create snapshot_cleaner --display-name "Snapshot Cleaner"
gcloud projects add-iam-policy-binding ${PROJECT_ID} --member serviceAccount:snapshot-cleaner@${PROJECT_ID}.iam.gserviceaccount.com --role projects/ed-aem6/roles/snapshot_cleaner
gcloud iam service-accounts keys create ~/.snapshot-cleaner/key.json --iam-account snapshot-cleaner@${PROJECT_ID}.iam.gserviceaccount.com

export GOOGLE_APPLICATION_CREDENTIALS=~/key.json
./google-cloud-snapshot-cleaner
