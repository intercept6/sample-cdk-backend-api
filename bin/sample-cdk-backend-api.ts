#!/usr/bin/env node
import 'source-map-support/register';
import {FrontendStack} from '../lib/frontend-stack';
import {SampleCdkBackendApiStack} from "../lib/sample-cdk-backend-api-stack";
import cdk = require('@aws-cdk/core');

const util = require('util');
const exec = util.promisify(require('child_process').exec);

async function deploy() {
    await exec('GO111MODULE=off go get -v -t -d ./src/backend/persons/add-person/... &&' +
        'GOOS=linux GOARCH=amd64 ' +
        'go build -o ./src/backend/persons/add-person/main ./src/backend/persons/add-person/**.go');
    await exec('GO111MODULE=off go get -v -t -d ./src/backend/persons/get-persons/... &&' +
        'GOOS=linux GOARCH=amd64 ' +
        'go build -o ./src/backend/persons/get-persons/main ./src/backend/persons/get-persons/**.go');


    const app = new cdk.App();
    new SampleCdkBackendApiStack(app, 'SampleCdkBackendApiStack');
    new FrontendStack(app, 'Frontend');

    // await exec('rm ./src/backend/persons/add-person/main');
    // await exec('rm ./src/backend/persons/get-persons/main');
}

deploy();
