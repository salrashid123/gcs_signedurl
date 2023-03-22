
import os

from flask import Flask

import google.auth
from google.auth import impersonated_credentials
from google.auth import compute_engine

from datetime import datetime, timedelta
from google.cloud import storage

app = Flask(__name__)


@app.route("/")
def hello_world():
    bucket_name = os.environ.get("BUCKET_NAME")
    sa_email =  os.environ.get("SA_EMAIL")


    credentials, project = google.auth.default()   

    storage_client = storage.Client()
    data_bucket = storage_client.bucket(bucket_name)
    blob = data_bucket.blob("file.txt")
    expires_at_ms = datetime.now() + timedelta(minutes=30)
    signing_credentials = impersonated_credentials.Credentials(
        source_credentials=credentials,
        target_principal=sa_email,
        target_scopes = 'https://www.googleapis.com/auth/devstorage.read_only',
        lifetime=2)

    # using the default credentials in the blob.generate_signed_url() will not work
    # signed_url = blob.generate_signed_url(expires_at_ms, credentials=credentials)
    signed_url = blob.generate_signed_url(expires_at_ms, credentials=signing_credentials)

    return signed_url, 200, {'Content-Type': 'text/plain'}


if __name__ == "__main__":
    app.run(debug=True, host="0.0.0.0", port=int(os.environ.get("PORT", 8080)))