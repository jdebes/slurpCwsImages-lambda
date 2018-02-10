# slurpCwsImages AWS Lambda Function

AWS Lambda that grabs codeswholesale images and writes them to an S3 bucket.

## Getting Started

This project manages dependencies with dep, use the following commands to build the vendor folder after doing a clone.

```
dep ensure
```

### Environment Variables

Codeswholesale:
```
CWS_CLIENT_ID
CWS_CLIENT_SECRET
CWS_CLIENT_TOKEN #https://api.codeswholesale.com/oauth/token
```

AWS:
```
AWS_ACCESS_KEY_ID #When running on local dev machine only
CWS_CLIENT_SECRET #When running on local dev machine only
BUCKET
```

### Run Tests

Run all tests:
```
go test ./...
```