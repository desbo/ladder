deploy: deploy-api deploy-static

deploy-api:
	gcloud app deploy api -q --project tt-ladder

deploy-static: 
	cd static && NODE_ENV=production yarn webpack
	gcloud app deploy static -q --project tt-ladder
