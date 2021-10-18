# s3-remove-exif
soon after uploaded, remove exif

## Installation

1. Set up AWS
   - IAM Policy
     - as below json
   - IAM Role
     - `AWSLambdaBasicExecutionRole`
     - the policy above
   - Lambda
     - Runtime: `Go 1.x`
     - Handler: `OnObjectCreated`
     - Arch: `x86_64`
     - IAM Role: the role above
   - S3
     - Event trigger: Lambda above
     - ACL, Bucket Policy: as you like
3. Download *s3-remove-exif.zip* from [Releases](https://github.com/usagiga/s3-remove-exif/releases/latest)
   - Or, run `curl -fsSL https://github.com/usagiga/s3-remove-exif/releases/download/${VERSION}/s3-remove-exif-lambda.zip --output s3-remove-exif-lambda.zip`
4. Upload *s3-remove-exif.zip* to AWS Lambda

```json:policy.json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "s3-remove-exif-default",
            "Effect": "Allow",
            "Action": [
                "s3:PutObject",
                "s3:GetObject",
                "s3:ListBucket",
            ],
            "Resource": [
                "arn:aws:s3:::YOUR_BUCKET_HERE/*",
                "arn:aws:s3:::YOUR_BUCKET_HERE"
            ]
        }
    ]
}
```

## Usage

1. Upload an image to your bucket
2. Soon after remove EXIF from uploaded image
   - NOTE: processed images has default ACL ( `private` )
3. :tada:

## LICENSE

MIT
