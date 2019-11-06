import * as Lambda from 'aws-lambda';

export const handler: Lambda.APIGatewayProxyHandler = async (proxyEvent: Lambda.APIGatewayEvent, _context) => {

    const petId = proxyEvent.body;
    console.log('test ');
    return {
        statusCode: 200,
        body: 'OK'
    }
};