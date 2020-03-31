# GCP Alert to Webhook

## Deploy Command

```sh
$ gcloud functions deploy AlertFunc --region {region} --env-vars-file .env.yaml --runtime go113 --trigger-http --allow-unauthenticated
```
