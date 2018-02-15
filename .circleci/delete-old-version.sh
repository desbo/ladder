set -o nounset
set -o errexit

gcloud auth activate-service-account --key-file ${HOME}/gcp-key.json

MAX_VERSIONS=15
NUM_VERSIONS=$(gcloud app versions list --project tt-ladder --format list | wc -l)

if [ $NUM_VERSIONS -ge $(($MAX_VERSIONS - 2)) ]; then
  gcloud app versions list \
    --project tt-ladder \
    --sort-by=LAST_DEPLOYED \
    --limit 1 \
    --service default \
    --format="value(version.id)" |
    xargs gcloud app versions delete --quiet --project tt-ladder
  
  gcloud app versions list \
    --project tt-ladder \
    --sort-by=LAST_DEPLOYED \
    --limit 1 \
    --service api \
    --format="value(version.id)" |
    xargs gcloud app versions delete --quiet --project tt-ladder

  exit 0
fi

echo "nothing to delete ($(expr $NUM_VERSIONS + 0) versions)"