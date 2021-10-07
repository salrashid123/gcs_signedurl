// DO NOT USE, google-auth impersonation does not work with storage

const { GoogleAuth, Impersonated } = require('google-auth-library');
const {Storage} = require('@google-cloud/storage');

const express = require('express');
const app = express();

app.get('/', async (req, res) => {

  const bucketName =  process.env.BUCKET_NAME;
  const saEmail =  process.env.SA_EMAIL;
  fileName = "file.txt"

  const auth = new GoogleAuth();
  const client = await auth.getClient();
  try {
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
    res.send(`Hello ${url}!`);
    }  catch (err) {
        console.log(err.stack);
        res.sendStatus(500);
    }
  
});

const port = process.env.PORT || 8080;
app.listen(port, () => {
  console.log(`helloworld: listening on port ${port}`);
});