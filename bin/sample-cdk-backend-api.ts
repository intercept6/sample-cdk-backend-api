#!/usr/bin/env node
import 'source-map-support/register';
import {FrontendStack} from '../lib/frontend-stack';
import cdk = require('@aws-cdk/core');

const app = new cdk.App();
// new SampleCdkBackendApiStack(app, 'SampleCdkBackendApiStack');
new FrontendStack(app, 'Frontend');