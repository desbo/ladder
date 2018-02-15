deploy: deploy-api deploy-static

deploy-api:
	gcloud app deploy api -q --project tt-ladder
	gcloud datastore create-indexes -q --project tt-ladder api/index.yaml

deploy-static: 
	cd static && NODE_ENV=production yarn webpack
	gcloud app deploy static -q --project tt-ladder
