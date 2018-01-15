deploy:
	cd static && NODE_ENV=production webpack
	gcloud app deploy api static --project tt-ladder