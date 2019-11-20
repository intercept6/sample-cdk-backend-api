import {Bucket} from '@aws-cdk/aws-s3';
import {CfnOutput, Construct, Duration, RemovalPolicy, Stack, StackProps} from "@aws-cdk/core";
import {BucketDeployment, Source} from '@aws-cdk/aws-s3-deployment';
import {CfnCloudFrontOriginAccessIdentity, CloudFrontWebDistribution} from '@aws-cdk/aws-cloudfront'
import {Effect, PolicyStatement} from '@aws-cdk/aws-iam';

export class FrontendStack extends Stack {
    constructor(scope: Construct, id: string, props?: StackProps) {
        super(scope, id, props);

        const websiteBucket = new Bucket(this, 'Website', {
            removalPolicy: RemovalPolicy.DESTROY
        });

        new BucketDeployment(this, 'DeployWebsite', {
            sources: [Source.asset('src/frontend/dist')],
            destinationBucket: websiteBucket,
        });

        // CloudFrontからwebsiteBucketにアクセスする際のOriginAccessIdentity
        const OAI = new CfnCloudFrontOriginAccessIdentity(this, 'OAI', {
            cloudFrontOriginAccessIdentityConfig: {
                comment: `WebsiteBucket-${this.stackName}`
            }
        });

        // webSiteBucketのBucketPolicyのStatement
        // 先ほど作ったOAIにs3:GetObjectを許可する
        // websiteBucketはpublic access出来ない設定（デフォルト）になっているので
        // こうしておかないとCloudFrontからアクセス出来ない
        const webSiteBucketPolicyStatement = new PolicyStatement({effect: Effect.ALLOW});
        webSiteBucketPolicyStatement.addCanonicalUserPrincipal(OAI.attrS3CanonicalUserId);
        webSiteBucketPolicyStatement.addActions("s3:GetObject");
        webSiteBucketPolicyStatement.addResources(`${websiteBucket.bucketArn}/*`);
        websiteBucket.addToResourcePolicy(webSiteBucketPolicyStatement);

        const distribution = new CloudFrontWebDistribution(this, 'WebsiteDistribution', {
            originConfigs: [
                {
                    s3OriginSource: {
                        s3BucketSource: websiteBucket,
                        originAccessIdentityId: OAI.ref
                    },
                    behaviors: [{
                        isDefaultBehavior: true,
                        minTtl: Duration.seconds(0),
                        maxTtl: Duration.seconds(0),
                        defaultTtl: Duration.seconds(0),
                    }]
                }
            ]
        });

        new CfnOutput(this, 'CFTopURL', {value: `https://${distribution.domainName}/`})
    }
}
