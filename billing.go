package main

import (
	"fmt"
	"strings"
	

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	// Configure your AWS session
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"), // Replace with your AWS region
	}))

	// Specify the S3 bucket and prefix where billing reports are stored
	bucketName := "your-billing-bucket"
	prefix := "aws-cost-explorer/AWSUsageReport/"

	// Create an S3 client
	svc := s3.New(sess)

	// List objects in the S3 bucket
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		fmt.Println("Error listing objects:", err)
		return
	}

	// Iterate through the billing report files
	for _, obj := range resp.Contents {
		if strings.HasSuffix(*obj.Key, ".csv") {
			// Download and process CSV billing report
			reportData, err := downloadS3Object(svc, bucketName, *obj.Key)
			if err != nil {
				fmt.Println("Error downloading report:", err)
				continue
			}

			// Process the billing report data
			// You will need to implement your logic to parse and extract billing information from the CSV data
			// For example, you can use a CSV parsing library like "encoding/csv" to parse the data.

			// Print the content of the report (for demonstration purposes)
			fmt.Println("Billing Report Content:")
			fmt.Println(reportData)
		}
	}
}

// downloadS3Object downloads an object from an S3 bucket and returns its content as a string
func downloadS3Object(svc *s3.S3, bucketName, objectKey string) (string, error) {
	resp, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the object content
	buf := make([]byte, 1024)
	var content []byte
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			content = append(content, buf[:n]...)
		}
		if err != nil {
			break
		}
	}

	return string(content), nil
}

