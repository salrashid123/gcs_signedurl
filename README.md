## Google Cloud Storage SignedURL with Cloud Run, Cloud Functions and GCE VMs

Code snippet to create a GCS Signed URL in Cloud Run, Cloud Functions and GCE VMs

- [Signed URL](https://cloud.google.com/storage/docs/access-control/signed-urls)
- [Cloud Run](https://cloud.google.com/run/docs)

Why am i writing this repo?  because it isn't clear that in those environment that you need to use [service account impersonation](https://cloud.google.com/iam/docs/impersonating-service-accounts) and [iamcredentials.serviceAccounts.signBlob](https://cloud.google.com/iam/docs/reference/credentials/rest/v1/projects.serviceAccounts/signBlob) 


### Setup

```bash
export PROJECT_ID=`gcloud config get-value core/project`
export PROJECT_NUMBER=`gcloud projects describe $PROJECT_ID --format="value(projectNumber)"`
export BUCKET_NAME=crdemo-$PROJECT_NUMBER

gsutil mb gs://$BUCKET_NAME

echo foo > file.txt
gsutil cp file.txt gs://$BUCKET_NAME

# allow cloud run's default service account access
gsutil acl ch -u $PROJECT_NUMBER-compute@developer.gserviceaccount.com:R gs://$BUCKET_NAME/file.txt

gcloud iam service-accounts  add-iam-policy-binding   --role=roles/iam.serviceAccountTokenCreator  \
 --member=serviceAccount:$PROJECT_NUMBER-compute@developer.gserviceaccount.com $PROJECT_NUMBER-compute@developer.gserviceaccount.com

gcloud config set run/region us-central1
```

Then build and push the language your'e interested in below

```bash
docker build -t gcr.io/$PROJECT_ID/crsigner .
docker push gcr.io/$PROJECT_ID/crsigner

gcloud run deploy signertest --image gcr.io/$PROJECT_ID/crsigner --platform=managed --set-env-vars="BUCKET_NAME=$BUCKET_NAME,SA_EMAIL=$PROJECT_NUMBER-compute@developer.gserviceaccount.com"

export CR_URL=`gcloud run services describe  signertest --format="value(status.url)"`

curl -s $CR_URL

curl -s `curl -s $CR_URL`
```

### golang

In golang, we're using the IAMCredentials api to sign the bytes.  [PR 4604](https://github.com/googleapis/google-cloud-go/pull/4604) seeks to automate that

### java

for java, first build
```
mvn clean install the image
```

then build the docker image, push to gcr then deploy to cloud run

### Python

[google-auth python](https://google-auth.readthedocs.io/en/master/) offers two ways to `signer` interfaces you can use:

* [compute_engine.IDTokenCredentials.signer](https://google-auth.readthedocs.io/en/master/reference/google.auth.compute_engine.html#google.auth.compute_engine.IDTokenCredentials.signer)

* [impersonated_credentials.Credentials.signer](https://google-auth.readthedocs.io/en/master/reference/google.auth.impersonated_credentials.html#google.auth.impersonated_credentials.Credentials.signer)

You might be wondering why an IDToken credentials has a signer?  Well, thats a side effect of an incorrect initial implementation of the compute engine ID Token (see [issue #344](https://github.com/googleapis/google-auth-library-python/issues/344)).   The interface users should use is `impersonated_credential`


### nodeJS

google-cloud node js auth library does not support signing using impersoantion
- see [issue #1210](https://github.com/googleapis/google-auth-library-nodejs/issues/1210)

Which means you need to do this by hand by using an authorized client

see [google-auth-library-nodejs#impersonated-credentials-client](https://github.com/googleapis/google-auth-library-nodejs#impersonated-credentials-client)

```javascript
    // https://cloud.google.com/iam/docs/reference/credentials/rest/v1/projects.serviceAccounts/signBlob
    const name = "projects/-/serviceAccounts/" + saEmail
    const url = 'https://iamcredentials.googleapis.com/v1/' + name +':signBlob'
    // construct the POST parameters with the bytes to sign and use targetClient.request
```

### dotnet

(contributions welcome)