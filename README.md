[![Build Status](https://travis-ci.org/kaihendry/goserverless.sg.svg?branch=master)](https://travis-ci.org/kaihendry/goserverless.sg)

Simple landing page for goserverless.sg which advertises Kai Hendry's Singapore
based serverless services and expertise.

`/rank` offers an API to show the number of services per region, which currently is fixed on the AWS SDK parameters. Hence the [daily deployment of the service](https://travis-ci.org/kaihendry/goserverless.sg/), since it should build will the latest SDK.

A better approach could be this API: https://twitter.com/hichaelmart/status/1054561536121937920

https://github.com/aws/aws-sdk-go-v2/issues/99
