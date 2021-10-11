const {Storage} = require('@google-cloud/storage');

const express = require('express');
const app = express();

app.get('/', async (req, res) => {

  const bucketName =  process.env.BUCKET_NAME;
  const saEmail =  process.env.SA_EMAIL;
  fileName = "file.txt"
  try {
    const storage = new Storage();

    const options = {
        version: 'v4',
        action: 'read',
        expires: Date.now() + 10 * 60 * 1000, // 10 minutes
    };

    const [url] = await storage
        .bucket(bucketName)
        .file(fileName)
        .getSignedUrl(options);
    
    res.send(`${url}!`);
    }  catch (err) {
        console.log(err.stack);
        res.sendStatus(500);
    }
  
});

const port = process.env.PORT || 8080;
app.listen(port, () => {
  console.log(`helloworld: listening on port ${port}`);
});