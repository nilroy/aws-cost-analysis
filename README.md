<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [aws-cost-analysis](#aws-cost-analysis)
  - [Generating report for current EC2 usage](#generating-report-for-current-ec2-usage)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# aws-cost-analysis
Contains tools to generate various reports to help analyize aws costs

## Generating report for current EC2 usage

Build the binary using below command

```
$ cd ec2-cost-analysis
$ make build
```

The binary is placed inside `ec2-cost-analysis/bin`

Run the binary

```
bin/ec2-cost-analysis -region=us-east-1 -s3 <s3 bucket name>
```

This would generate a CSV file that contains the Instances that we are curently running in all environments and associated cost. The program automatically uploads the CSV to the S3 bucket. We are using Amazon QuickSight that ingests that data and updates the dashboards that shows comparison between different pricing models for the current EC2 instane usage.