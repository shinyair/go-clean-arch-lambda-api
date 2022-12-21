# go-clean-arch project with AWS Lambda and AWS API Gateway
This is a serverless go program which implements the clean architecture and integrates with AWS Lambda, API Gateway, and DynamoDB.

## Prerequisites
#### npm
npm is required by serverless framework, and also we can use scripts in package.json to build & deploy our service.
[Downloading and installing Node.js and npm](https://docs.npmjs.com/downloading-and-installing-node-js-and-npm)

#### aws-cli
aws-cli is required to debug your program in local and deploy your personal online env from local.
[Getting started with the AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html)

#### go
[Tutorial: Get started with Go](https://go.dev/doc/tutorial/getting-started)

#### vscode
It's FREE!

**`Required plugins:`**
- Go, which is provided by Go Team at Google. It is required to develop go programs in vscode

**`Suggested plugins:`**
- GoComment, which is provided by lijunawk. It helps generate comments for functions and structs
- docs-yaml, which is provided by Microsoft. It helps format serverless.yml in vscode

#### serverless
Serverless Framework is easy for starup to understand the architecture of services and also it support not only AWS but also Google Cloud Platform and Microsoft Azure. Instal Serverless Framework by npm: `npm install -g serverless`.
[Serverless.yml Reference](https://www.serverless.com/framework/docs/providers/aws/guide/serverless.yml)

#### 7z(for windows only)
Windows desn't support zip files by cmd natively, so use 7z cmd line to support package zip for deployment by scripts.
[Download 7z](https://www.7-zip.org/download.html)
[Unzip files (7-zip) via cmd command - stackoverflow](https://stackoverflow.com/questions/14122732/unzip-files-7-zip-via-cmd-command)

## Project structure
Consider implementing `golang-standard` and the `clean architecture` in lambda based progoram (which always serves functions), the project is organized as:
```
project 
├── cmd 
|    ├── main.go # lambda handler, use mux proxy provide by aws
|    └── local
|         └── main.go # help run & debug local
├── configs 
|    └── env.yaml # env variable config file for all stage
├── deployment 
|    ├── db
|    |    └── serverless.yml # serverless to deploy db. need to be deployed before
|    ├── serverless.yml # serverless to deploy lambda, apigateway and so on
|    ├── output # folder generated by build script in package.json. gitignore it
|    └── output.zip # file generated by package script in package.json. gitignore it
├── internal # clean architecture
|    ├── app
|    |    ├── app.go # init beans
|    |    └── app_config.go # handle configs
|    ├── controller # handle http requests by mux
|    |    ├── biz
|    |    |    └── dummy_controller.go # entrance of dummy logic
|    |    ├── controller.go # interface
|    |    ├── mux_controller.go # mux implementatation of controller.go
|    |    ├── mux_middlewares.go # implement necessary mux middlewares(also known as interseptors/filters)
|    |    └── mux_router.go # helper function to init mux router
|    ├── domain
|    |    └── dummy.go # entity and repository interface
|    ├── repository
|    |    └── dynamodb
|    |         └── dummy_dynamodb_repo.go # dynamodb implementation of dummy repository
|    ├── usecase
|    |    └── dummy_usecase # biz logic of dummy
|    └── logger
|         └── logger.go # logger utils
├── scripts
├── .gitignore
├── go.mod
├── go.sum
├── package.json
├── package-lock.json
└── README.md
```
## Build the existing project
Run cmd `go mod tidy` to install all required packages.

## Test in local
### Check env.yml

**`Required`**
- change the `VARIANT`. `VARIANT` helps deploy your own aws services online for each stage. For example, name of official lambda of `dev` stage begins as `dev--{appcode}-`, while the name of your personal lambda of `dev` stage begins as `dev-{variant}-{appcode}-`.
- input the `ACCOUNT_ID`. ID of AWS Account you are using.
- change the `AWS_PROFILE` if you have a specific profile. remove it if you do not have the specific profile.

**`Optional`**
- change the `APPCODE` as you want.
- change the `REGION` and `AWS_REGION` to the one you want deploy resources to. 
- change the `AWS_DEPLOYMENT_BUCKET` as you want.
- change the `DUMMY_TABLE_NAME` in format `{stage}-{appcode}-dummy`.

### Before testing
When debuging from local without localstack, aws services should be prepared in advance except API Gateway & Lambda, such as DynamoDB, so we need to deploy DynamoDB by Serverless Framework before testing.
[Reason to create dynamodb tables as a service](https://stackoverflow.com/questions/41620437/how-to-continue-deploy-if-dynamodb-table-already-exists)

- Add a DynamoDB Serverless service
  - open folder: `go-clean-arch-lambda-api/deployment/db`
  - check [serverless example file](https://github.com/serverless/examples/blob/v3/aws-golang-rest-api-with-dynamodb/serverless.yml)
  - add a serverless service for dynamodb
  - use variables defined in `go-clean-arch-lambda-api/configs` to create dynamodb serverless resource. check `deployment/db/serverless.yml` for detail
- Deploy the created DynamoDB Serverless service
  - open folder: `go-clean-arch-lambda-api/deployment/db`
  - run cmd `npm install serverless-deployment-bucket --save-dev`
  - run cmd `sls deploy`
  - if you have a aws profile, use `--aws-profile {your profile name}` cmd param

### Run local main and test
- run cmd `go run ./cmd/local/main.go` under `go-clean-arch-lambda-api`
- test `get`/`post`/`delete` dummy api by Postman or the other tools

## Deployment
### Build and package by npm
- For Windows, we need to set go env `GOARCH` as `amd64`, `GOOS` as `linux`, because the online lambda runtime is based on linux core. check go environment variables by `go env`. 
  - [How to cross compile from Windows to Linux - stackoverflow](https://stackoverflow.com/questions/20829155/how-to-cross-compile-from-windows-to-linux)
  - [How to use environment variables in NPM](https://blog.jimmydc.com/cross-env-for-environment-variables/)
- open folder: `go-clean-arch-lambda-api`
  - run cmd `npm install copyfiles --save-dev`
  - run cmd `npm install rimraf --save-dev`
  - add following scripts in of package.json
  ```
  "scripts": {
    "prebuild": "rm -rf deployment/output/* && copyfiles --flat configs/* deployment/output/configs",
    "build": "set GOARCH=amd64&& set GOOS=linux&& go build -o deployment/output/bin/main cmd/main.go",
    "test": "echo \"Error: no test specified\" && exit 1"
  },
  ```
  - run cmd `npm run build`
  - run cmd `npm run package`
  - check `deployment/outout` folder, the `configs` folder is copied & pasted, and a `main` file is generated


### Deploy Serverless services
#### Deploy DynamoDB service
refs `Test in local#Before testing`

#### Deploy Lambda and API Gateway
- serverless framework read the aws account and your aws credentials from your local credentials file, so make sure your credentials are valid. [Configuration and credential file settings](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html)
- open folder: `go-clean-arch-lambda-api`
  - if you want use a specific local AWS profile to deploy, add the param ` --aws-profile {profile_name}` in deploy scripts of the `package.json`, like
   ```
  "scripts": {
    ...,
    "deploy": "cd deployment && sls deploy --force --aws-profile profileA",
    "undeploy": "cd deployment && sls remove --aws-profile profileA",
    ...
  },
  ```
  - run cmd `npm run deploy`

## Test on AWS
### Run and check log
**`Invoke URL`**
Open API Gateway service on AWS Console, and find the deployed API Gateway by region and name `{stage}-{variant}-go-clean-lambda` (which are set in serverless.yml).
Find url of your API Gateway, and try to invoke:
- run cmd `npm install curl -g`
- run cmd `curl {api_gateway_invoke_url}/api/dummy/1` and check the result
- run cmd `curl -X POST {api_gateway_invoke_url}/api/dummy??id=1&name=aaa&attr=ttt`
- run cmd `curl {api_gateway_invoke_url}/api/dummy/1` and check the result
- run cmd `curl -X DELETE {api_gateway_invoke_url}/api/dummy/1`
- run cmd `curl {api_gateway_invoke_url}/api/dummy/1` and check the result

**`Check Lambda Log`**
Open Lambda service on AWS Console, and find the deployed Lambda function by region and name `{stage}-{variant}-go-clean-lambda-dummy` (which are set in serverless.yml).
Switch to `Monitor` tab and open `Cloudwatch` to check the runtime log.

### Invoked AWS services
- API Gateway, trigger of Lambda
- Lambda, function to run the go program
- DynamoDB, db
- Cloudwatch, logs
- IAM, policies and roles to run Lambda
- S3, deployment bucket, where zips are uploaded
- CloudFormation, deployment stack, which manages deployed resources of each serverless service. Here are 2 stacks in the project, one for db, one for lambda-api.

## Reference
### Standard Go Project Layout
[Standard Go Project Layout](https://github.com/golang-standards/project-layout)
```
project
├── api
├── assets
├── build
├── cmd
├── configs
├── deployments
├── docs
├── examples
├── githooks
├── init
├── internal
├── pkg
├── scripts
├── test
├── third_party
├── tools
├── vendor
├── web
└── website
```
### The Clean Code Blog
[The Clean Code Blog](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
<img src="https://blog.cleancoder.com/uncle-bob/images/2012-08-13-the-clean-architecture/CleanArchitecture.jpg">

## Appendix
### Create new projects 
#### Create the npm project
Run cmd `npm init` to create the `package.json`. For example:
```
> npm init
package name: (go-clean-arch-lambda-api)
version: (1.0.0) 0.0.1
description:
entry point: (index.js)
test command:
git repository:
keywords:
author:
license: (ISC)
```

#### Create the go project
Run cmd to init a new module, such as `go mod init local.com/go-clean-lambda`.
Create `golang-standard` structure like:
```
project  
├── cmd 
├── configs 
├── deployment 
├── internal 
├── scripts
└── README.md
```
Write internal code follow the `clean architecture`. Example:
```
project 
├── cmd 
|    └── main.go # lambda handler
├── configs
├── deployment # serverless.yml
├── internal # clean architecture
|    ├── app
|    ├── controller
|    ├── domain
|    ├── repository
|    |    └── dynamodb
|    ├── usecase
|    └── utils
├── scripts
├── .gitignore
├── go.mod
├── go.sum
├── package.json
├── package-lock.json
└── README.md
```
Write your own code in each layter and run cmd `go get xxxx` to install a new package when necessary.

#### Add scripts to build, packge, and deploy the program
For Windows, we need to set go env `GOARCH` as `amd64`, `GOOS` as `linux`, because the online lambda runtime is based on linux core. check go environment variables by `go env`. 
- [How to cross compile from Windows to Linux - stackoverflow](https://stackoverflow.com/questions/20829155/how-to-cross-compile-from-windows-to-linux)
- [How to use environment variables in NPM](https://blog.jimmydc.com/cross-env-for-environment-variables/)

Add scripts in `package.json` like:
```
  "scripts": {
    "prebuild": "rm -rf deployment/output/* && copyfiles --flat configs/* deployment/output/configs",
    "build": "set GOARCH=amd64&& set GOOS=linux&& go build -o deployment/output/bin/main cmd/main.go",
    "deploy": "cd deployment && sls deploy --force",
    "undeploy": "cd deployment && sls remove",
    "test": "echo \"Error: no test specified\" && exit 1"
  },
```