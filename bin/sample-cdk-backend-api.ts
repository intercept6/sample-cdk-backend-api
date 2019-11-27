#!/usr/bin/env node
import 'source-map-support/register';
import {FrontendStack} from '../lib/frontend-stack';
import {BackendStack} from "../lib/backend-stack";
import cdk = require('@aws-cdk/core');
import AWS = require ('aws-sdk');
import fs = require('fs');


const util = require('util');
const exec = util.promisify(require('child_process').exec);
const writeFile = util.promisify(fs.writeFile);

const lambda_func_basedir = '/src/backend/persons';
const lambda_func_dir = [`${lambda_func_basedir}/add-person`, `${lambda_func_basedir}/get-persons`,
    `${lambda_func_basedir}/del-person`];

async function deploy() {
    const app = new cdk.App();

    for (let dir of lambda_func_dir) {
        await exec(`GO111MODULE=off go get -v -t -d .${dir}/... && ` +
            'GOOS=linux GOARCH=amd64 ' +
            `go build -o .${dir}/main .${dir}/**.go`);
    }
    const backend = await new BackendStack(app, 'BackendStack');
    const cfn = new AWS.CloudFormation();
    const backendStack = await cfn.describeStacks({StackName: 'BackendStack'}).promise();

    // CFnOutputsからAPI URLを取得する
    let apiUrl: string = '';
    if (
        backendStack.Stacks != undefined &&
        backendStack.Stacks[0].Outputs != undefined
    ) {
        const ApiUrlOutput = backendStack.Stacks[0].Outputs.find(item => item.OutputKey === 'ApiUrl');
        if (ApiUrlOutput != undefined) {
            apiUrl = ApiUrlOutput.OutputValue as string;
        }
    }
    // Envファイルに出力
    const data = 'NODE_ENV=\'production\'\n' +
        `VUE_APP_API_BASE_URL=\'${apiUrl}\'\n`;
    await writeFile('./src/frontend/.env', data);

    await exec('cd ./src/frontend/ && npm run build --fix');
    const frontend = new FrontendStack(app, 'FrontendStack');
    frontend.addDependency(backend);

    app.synth();

    for (let dir of lambda_func_dir) {
        await exec(`rm .${dir}/main`);
    }
}

deploy();
