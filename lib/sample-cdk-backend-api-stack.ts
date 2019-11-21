import cdk = require('@aws-cdk/core');
import {Code} from '@aws-cdk/aws-lambda';
import {LambdaBackend} from './lambda-backend';
import {Cors, RestApi} from '@aws-cdk/aws-apigateway';
import {AttributeType, BillingMode, Table} from '@aws-cdk/aws-dynamodb';


export class SampleCdkBackendApiStack extends cdk.Stack {
    constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
        super(scope, id, props);

        // API Gateway
        const api = new RestApi(this, 'RestApi', {
            restApiName: 'SampleCdkBackendApi',
            defaultCorsPreflightOptions: {
                allowOrigins: Cors.ALL_ORIGINS,
                allowCredentials: true,
                allowMethods: Cors.ALL_METHODS,
            }
        });

        // DynamoDB
        const personsTable = new Table(this, 'PersonsTable', {
            billingMode: BillingMode.PAY_PER_REQUEST,
            partitionKey: {name: 'Id', type: AttributeType.STRING}
        });

        // Lambda Functions
        // /persons
        const personsPath = api.root.addResource('persons');

        const getPersonsFunc = new LambdaBackend(this, 'GetPersons', {
            code: Code.fromAsset('./src/backend/persons/get-persons'),
            resource: personsPath,
            method: 'GET',
            environment: {
                'TABLE_NAME': personsTable.tableName
            }
        });
        personsTable.grantReadData(getPersonsFunc);

        const addPersonFunc = new LambdaBackend(this, 'AddPerson', {
            code: Code.fromAsset('./src/backend/persons/add-person'),
            resource: personsPath,
            method: 'POST',
            environment: {
                'TABLE_NAME': personsTable.tableName
            }
        });
        personsTable.grantReadWriteData(addPersonFunc);
    }
}
