## Google Cloud Storage SignedURL with Cloud Run, Cloud Functions and GCE VMs

Code snippet to create a GCS Signed URL in Cloud Run, Cloud Functions and GCE VMs

- [Signed URL](https://cloud.google.com/storage/docs/access-control/signed-urls)
- [Cloud Run](https://cloud.google.com/run/docs)

- Why am i writing this repo?  
  because it isn't clear that in those environment that with _some languages_ you can "just use" the default credentials (`node`, `java`) while in others you need to explicitly  use [service account impersonation](https://cloud.google.com/iam/docs/impersonating-service-accounts) (`go`)...and finally in `python`, you _should_ use the the impersonated credential type directly to sign.  Documented [samples for signedURL](https://cloud.google.com/storage/docs/samples/storage-generate-signed-url-v4) uses actual keys!  don't do this!

- Why Impersonation?  
  Well, Cloud Run, Cloud Functions and GCE environments do not have anyway to sign anything (and no, do NOT embed a service account key file anywhere!).  Since those environments can sign by themselves, they need to use an API to sign on behalf of itself.  That API is is listed above

- Why are they different in different languages?  
  `java`, `node` automatically detects that its running in a Cloud Run|GCE and "knows" it can't sign by itself and instead attempts to use impersonation automatically.  In `go`, it doesn't do that so you need to directly use that API.  In `python`, all the examples i saw around incorrectly uses the wrong Credential type to sign.

Impersonated signing uses ([iamcredentials.serviceAccounts.signBlob](https://cloud.google.com/iam/docs/reference/credentials/rest/v1/projects.serviceAccounts/signBlob).

The following examples uses the ambient service account to sign (eg, the service account cloud run uses is what signs the URL).  It would be a couple more steps to make cloud run sign on behalf of _another_ service account (and significantly more steps for java and node that made some assumptions on your behalf already)

Finally, if you must use a key, try to embed it into hardware, if possible.
  - [GCS signedURLs and GCP Authentication with Trusted Platform Module](https://medium.com/google-cloud/gcs-signedurls-and-gcp-authentication-with-trusted-platform-module-482faff2ac04)
### Setup

```bash
export PROJECT_ID=`gcloud config get-value core/project`
export PROJECT_NUMBER=`gcloud projects describe $PROJECT_ID --format="value(projectNumber)"`
export BUCKET_NAME=crdemo-$PROJECT_NUMBER
export SA_EMAIL=$PROJECT_NUMBER-compute@developer.gserviceaccount.com

gsutil mb gs://$BUCKET_NAME

echo foo > file.txt
gsutil cp file.txt gs://$BUCKET_NAME

# allow cloud run's default service account access
gsutil acl ch -u $SA_EMAIL:R gs://$BUCKET_NAME/file.txt

gcloud iam service-accounts  add-iam-policy-binding   --role=roles/iam.serviceAccountTokenCreator  \
 --member=serviceAccount:$SA_EMAIL $SA_EMAIL

gcloud config set run/region us-central1
```

Then build and push the language your'e interested in below

```bash
docker build -t gcr.io/$PROJECT_ID/crsigner .
docker push gcr.io/$PROJECT_ID/crsigner

gcloud run deploy signertest --image gcr.io/$PROJECT_ID/crsigner --platform=managed --set-env-vars="BUCKET_NAME=$BUCKET_NAME,SA_EMAIL=$SA_EMAIL"

export CR_URL=`gcloud run services describe  signertest --format="value(status.url)"`

curl -s $CR_URL

curl -s `curl -s $CR_URL`
```

### golang

In golang, we're using the IAMCredentials api to sign the bytes.  [PR 4604](https://github.com/googleapis/google-cloud-go/pull/4604) seeks to automate that

### java

for java, first build
```
mvn clean install
```

then build the docker image, push to gcr then deploy to cloud run

### Python

[google-auth python](https://google-auth.readthedocs.io/en/master/) offers two ways to `signer` interfaces you can use:

* [compute_engine.IDTokenCredentials.signer](https://google-auth.readthedocs.io/en/master/reference/google.auth.compute_engine.html#google.auth.compute_engine.IDTokenCredentials.signer)

* [impersonated_credentials.Credentials.signer](https://google-auth.readthedocs.io/en/master/reference/google.auth.impersonated_credentials.html#google.auth.impersonated_credentials.Credentials.signer)

You might be wondering why an IDToken credentials has a signer?  Well, thats a side effect of an incorrect initial implementation of the compute engine ID Token (see [issue #344](https://github.com/googleapis/google-auth-library-python/issues/344)).   The interface users should use is `impersonated_credential`


### nodeJS

Nodejs its pretty easy since by default, the library _automatically_ tries to use IAMCredentials API in these environments;

see [bucket.getSignedUrl()](https://googleapis.dev/nodejs/storage/latest/Bucket.html#getSignedUrl)

```
In Google Cloud Platform environments, such as Cloud Functions and App Engine, 
you usually don't provide a keyFilename or credentials during instantiation. In those environments, we call the signBlob API
```

However, my preference would've been to make it explicit applied which would also allow you to set a different account to sign with.  For example, like in the sample below using the IAM [google-auth-library-nodejs#impersonated-credentials-client](https://github.com/googleapis/google-auth-library-nodejs#impersonated-credentials-client)

```javascript
const { GoogleAuth, Impersonated } = require('google-auth-library');
const {Storage} = require('@google-cloud/storage');
const {IAMCredentialsClient} = require('@google-cloud/iam-credentials');

    let targetClient = new Impersonated({
        sourceClient: client,
        targetPrincipal: saEmail,
        lifetime: 10,
        delegates: [],
        targetScopes: ['https://www.googleapis.com/auth/cloud-platform']
    });

    const storage = new Storage({
        auth: {
            getClient: () => targetClient,
        },
    });

    const options = {
        version: 'v4',
        action: 'read',
        expires: Date.now() + 10 * 60 * 1000, // 10 minutes
    };

    const [url] = await storage
        .bucket(bucketName)
        .file(fileName)
        .getSignedUrl(options);
```

Unfortunately, google-cloud node js auth storage library does not support signing using impersonation with signer
- see [issue #1210](https://github.com/googleapis/google-auth-library-nodejs/issues/1210)


### dotnet

(contributions welcome)