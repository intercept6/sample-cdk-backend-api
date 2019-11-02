#!/usr/bin/env node
import 'source-map-support/register';
import {SampleCdkBackendApiStack} from '../lib/sample-cdk-backend-api-stack';
import cdk = require('@aws-cdk/core');

const app = new cdk.App();
new SampleCdkBackendApiStack(app, 'SampleCdkBackendApiStack');
