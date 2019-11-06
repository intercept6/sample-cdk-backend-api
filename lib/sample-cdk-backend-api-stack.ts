import cdk = require('@aws-cdk/core');
import {Code} from '@aws-cdk/aws-lambda';
import {LambdaBackend} from './lambda-backend';
import {RestApi} from '@aws-cdk/aws-apigateway';
import {AttributeType, BillingMode, Table} from '@aws-cdk/aws-dynamodb';


export class SampleCdkBackendApiStack extends cdk.Stack {
    constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
        super(scope, id, props);

        const api = new RestApi(this, 'RestApi', {
            restApiName: 'SampleCdkBackendApi',
        });

        // DynamoDB
        const petsTable = new Table(this, 'PetsTable', {
            billingMode: BillingMode.PAY_PER_REQUEST,
            partitionKey: {name: 'id', type: AttributeType.STRING}
        });

        // pets
        // /pets/{pet_id}
        const petsPath = api.root.addResource('pets');
        const petIdPath = petsPath.addResource('petId');

        // const getPetFunc = new LambdaBackend(this, 'GetPet', {
        //     code: Code.fromAsset('./src'),
        //     resource: petIdPath,
        //     method: 'GET'
        // });
        // petsTable.grantReadData(getPetFunc);

        const addPetFunc = new LambdaBackend(this, 'AddPet', {
            code: Code.fromAsset('./src/backend/' +
                'pets/add-pet'),
            resource: petsPath,
            method: 'POST'
        });
        petsTable.grantWriteData(addPetFunc);
    }
}
