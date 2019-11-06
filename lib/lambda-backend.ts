import {Construct} from '@aws-cdk/core';
import {Code, Function, Runtime} from '@aws-cdk/aws-lambda';
import {RetentionDays} from '@aws-cdk/aws-logs';
import {LambdaIntegration, Resource} from '@aws-cdk/aws-apigateway';

export type MethodType = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'

export interface LambdaFunctionProps {
    runtime?: Runtime
    handler?: string
    code: Code
    logRetention?: RetentionDays
    resource: Resource
    method: MethodType
}

export class LambdaBackend extends Function {
    constructor(scope: Construct, id: string, props: LambdaFunctionProps) {

        super(scope, id, {
            runtime: props.runtime || Runtime.NODEJS_10_X,
            handler: props.handler || 'index.handler',
            code: props.code,
            logRetention: props.logRetention || RetentionDays.TWO_WEEKS
        });

        const integration = new LambdaIntegration(this);
        props.resource.addMethod(props.method, integration)
    }
}
