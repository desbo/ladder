# table tennis ladder app

## setup
1. install app tools: https://cloud.google.com/sdk/docs/
2. `cd static && yarn install`

## running locally
1. run the API (http://localhost:8080): `dev_appserver.py --port 8080 --default_gcs_bucket_name tt-ladder.appspot.com api`
2. run the front-end (http://localhost:8081): `cd static && yarn dev`

there are also VSCode tasks set up for the above

## some likely queries
- create player (user)
- create ladder
- add player to ladder
- submit game
- get ladder
- get all ladders for a player (ones they own, and ones they're a part of)
- get game
- get player (including games played)