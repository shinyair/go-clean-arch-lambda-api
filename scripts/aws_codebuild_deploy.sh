#!/bin/sh
# how to create a pipeline for cicd:
# 1. create a full access role for codebuild (TODO: change to limited access role)
#  - trusted entity type: aws service
#  - use case: aws codebuild service
# 2. create a codebuild project for build & deploy serverless
#  - project name: xxx
#  - environment image: managed image
#  - operation system: ubuntu
#  - runtime: standard
#  - image: aws/codebuild/standard6.0
#  - service role: existing role, copy arn of the role created in 1st
#  - configure the others accordingly
# 3. open aws code pipeline service in aws management console and create a new pipeline
#  - pipeline settings stage
#    - pipeline name: xxx
#    - service role: new service role
#  - source stage
#    - source provider: github version2
#    - connection: select github or connect to github
#    - select repo and branch
#  - build stage
#    - build provider: codebuild
#    - project name: select codebuild project in 2nd step
#  - deploy stage: skip (build & deploy in codebuild directly)

# clear deployment of service
clear(){
    echo ">> clear"
    local SERVICE=$1
    if [ -z "$SERVICE" ]; then
        npm run clear
    else
        npm run clear:${SERVICE}
    fi
    echo "<< done"
}

# deploy db service
deploy_db() {
    echo ">> deploy dynamodb"
    npm run copy:configs:db
    npm run deploy:db
    echo "<< done"
}

# deploy all lambda services
deploy_lambda() {
    echo ">> build & deploy lambda"
    npm run test
    npm run lint
    npm run jwt # generate random jwt key pair config
    npm run copy:configs:lambda # copy configs to deploytment
    npm run deploy:lambda
    echo "<< done"
}

# deploy api doc
deploy_apidoc() {
    # deploy swagger
    echo ">> deploy swagger"
    npm run swag
    npm run copy:configs:apidoc
    npm run copy:apidoc
    npm run deploy:apidoc
    echo "<< done"
}

deploy() {
    local SERVICE=$1
    clear "$SERVICE"
    if [ -z "$SERVICE" ]; then
        echo "deploy all"
        deploy_db
        deploy_lambda
        deploy_apidoc
    elif [ "$SERVICE" == "db" ]; then
        deploy_db
    elif [ "$SERVICE" == "lambda" ]; then
        deploy_lambda
    elif [ "$SERVICE" == "apidoc" ]; then
        deploy_apidoc
    else
        ehco "unknown service ${SERVICE}"
    fi
}

deploy "$@"