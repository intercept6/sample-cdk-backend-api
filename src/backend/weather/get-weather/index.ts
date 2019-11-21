import * as Lambda from 'aws-lambda';
// import axios from 'axios';

export const handler: Lambda.APIGatewayProxyHandler = async (proxyEvent: Lambda.APIGatewayEvent, _context) => {

    return {
        statusCode: 200,
        headers: {
            "Access-Control-Allow-Origin": "*",
            "Access-Control-Allow-Credentials": "true"
        },
        body: '{"key": "value"}'
    }
};