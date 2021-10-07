package com.test;

import java.io.IOException;
import java.net.URL;
import java.util.Arrays;
import java.util.concurrent.TimeUnit;

import com.google.auth.oauth2.GoogleCredentials;
import com.google.auth.oauth2.ImpersonatedCredentials;
import com.google.cloud.storage.BlobId;
import com.google.cloud.storage.BlobInfo;
import com.google.cloud.storage.HttpMethod;
import com.google.cloud.storage.Storage;
import com.google.cloud.storage.Storage.SignUrlOption;
import com.google.cloud.storage.StorageOptions;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.http.MediaType;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.ResponseBody;


@SpringBootApplication
public class TestApp {

  @Controller
  class HelloworldController {
    @RequestMapping(value="/", produces = MediaType.TEXT_PLAIN_VALUE)
	@ResponseBody
    String hello() {

		String bucketName = System.getenv("BUCKET_NAME");
		String saEmail = System.getenv("SA_EMAIL");
		String objectName = "file.txt";


		String signedURLString = "";
		GoogleCredentials sourceCredentials;
		try {
			sourceCredentials =  GoogleCredentials.getApplicationDefault();
			ImpersonatedCredentials targetCredentials = ImpersonatedCredentials.create(sourceCredentials, saEmail, null,Arrays.asList("https://www.googleapis.com/auth/devstorage.read_only"), 2);

			Storage storage_service = StorageOptions.newBuilder().setCredentials(targetCredentials).build().getService();      
			
			BlobId blobId = BlobId.of(bucketName, objectName);
			BlobInfo blobInfo = BlobInfo.newBuilder(blobId).build();
			URL signedUrl = storage_service.signUrl(blobInfo, 600,  TimeUnit.SECONDS, 	SignUrlOption.httpMethod(HttpMethod.GET));
			signedURLString = signedUrl.toExternalForm();
	
		} catch (IOException ioex) {
			return "Error " + ioex.getMessage();
		}


      return signedURLString;
    }
  }

  public static void main(String[] args){
    SpringApplication.run(TestApp.class, args);
  }
}
