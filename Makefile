deploy:
	cd static && NODE_ENV=production webpack
	gcloud app deploy api static --project tt-ladder

deploy-api:
	gcloud app deploy api --project tt-ladder

deploy-static: 
	cd static && NODE_ENV=production webpack
	gcloud app deploy static --project tt-ladder
