import {expect as expectCDK, MatchStyle, matchTemplate} from '@aws-cdk/assert';
import cdk = require('@aws-cdk/core');
import SampleCdkBackendApi = require('../lib/sample-cdk-backend-api-stack');

test('Empty Stack', () => {
    const app = new cdk.App();
    // WHEN
    const stack = new SampleCdkBackendApi.BackendStack(app, 'MyTestStack');
    // THEN
    expectCDK(stack).to(matchTemplate({
        "Resources": {}
    }, MatchStyle.EXACT))
});