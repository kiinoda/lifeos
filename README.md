# LifeOS

Organize your life in a Google sheet, akin to a digital Bullet Journal with a weekly plan. If you're happy using a calendar application, this is not for you; at most, it might complement it.

Deploy lifeOS to:

* receive a daily schedule by email
* be notified of upcoming events
* receive notifications for monthly events

Tech-wise, this is an exercise in Go development, coupled with AWS Lambda, with a configuration kept in AWS SSM Parameter Store and the Serverless Framework.

Tech stack:

* Google sheet stores the data
* Go authenticates against Google Sheets API to retrieve data
* Credentials and configuration parameters are in AWS SSM Parameter Store
* Serverless Framework is used to deploy to AWS Lambda and set AWS IAM Permissions
* Emails are being sent using AWS SES
* All time events are being triggered using AWS EventBridge (CloudWatch Events)
