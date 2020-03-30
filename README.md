# GCP Alert to Webhook

## Deploy Command

```sh
$ gcloud functions deploy AlertFunc --region asia-northeast1 --env-vars-file .env.yaml --runtime go113 --trigger-http --allow-unauthenticated
```