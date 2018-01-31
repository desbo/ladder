set -o nounset
set -o errexit

gcloud auth activate-service-account --key-file ${HOME}/gcp-key.json

MAX_VERSIONS=15
NUM_VERSIONS=$(gcloud app versions list --project tt-ladder --format list | wc -l)

if [ $NUM_VERSIONS -ge $MAX_VERSIONS ]; then
  OLDEST_VERSION=$(gcloud app versions list \
    --project tt-ladder \
    --sort-by=LAST_DEPLOYED \
    --limit 1 \
    --format="value(version.id)")
 
  gcloud app versions delete --quiet --project tt-ladder $OLDEST_VERSION

  echo "deleted $OLDEST_VERSION"
fi

echo "nothing to delete ($(expr $NUM_VERSIONS + 0) versions)"