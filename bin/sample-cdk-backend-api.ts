#!/usr/bin/env node
import 'source-map-support/register';
import {FrontendStack} from '../lib/frontend-stack';
import {BackendStack} from "../lib/backend-stack";
import cdk = require('@aws-cdk/core');

const util = require('util');
const exec = util.promisify(require('child_process').exec);

const lambda_func_basedir = '/src/backend/persons';
const lambda_func_dir = [`${lambda_func_basedir}/add-person`, `${lambda_func_basedir}/get-persons`,
    `${lambda_func_basedir}/del-person`];

async function deploy() {
    for (let dir of lambda_func_dir) {
        await exec(`GO111MODULE=off go get -v -t -d .${dir}/... &&` +
            'GOOS=linux GOARCH=amd64 ' +
            `go build -o .${dir}/main .${dir}/**.go`);
    }

    const app = new cdk.App();
    new BackendStack(app, 'BackendStack');
    new FrontendStack(app, 'FrontendStack');
    app.synth();

    for (let dir of lambda_func_dir) {
        await exec(`rm .${dir}/main`);
    }
}

deploy();
