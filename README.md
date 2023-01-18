[![.github/workflows/sam.yml](https://github.com/kaihendry/goserverless.sg/actions/workflows/sam.yml/badge.svg)](https://github.com/kaihendry/goserverless.sg/actions/workflows/sam.yml)

Simple landing page for goserverless.sg which advertises Kai Hendry's Singapore
based serverless services and expertise.

`/rank` offers an API to show the number of services per region, which
currently is fixed on the AWS SDK parameters. Therefore it deploys daily.

Have to dig out Cloudfront CNAME from:
https://ap-southeast-1.console.aws.amazon.com/apigateway/main/publish/domain-names?domain=gosls.dabase.com
