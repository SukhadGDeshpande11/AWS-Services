package main

import (
    "fmt"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/ssm"
    "github.com/aws/aws-sdk-go/service/ec2"
    "github.com/aws/aws-sdk-go/service/cloudwatch"
    "github.com/aws/aws-sdk-go/service/cloudwatchlogs"
    "strings"
    "unicode"
    "time"
)

const (
    ParamPath = "/aws/service/global-infrastructure/regions/%s/services"
    // Add your S3 bucket name here
    S3BucketName = "your-s3-bucket-name"
)

func fetchEnabledServicesInRegion(region *string, sess *session.Session, cfg *aws.Config) (map[string]bool, error) {
    // ... your existing code for fetching AWS services ...

    return service, nil
}

func formatServiceName(s string) string {
    // ... your existing code for formatting service names ...
    return name
}

func main() {
    region := "ap-south-1"
    sess := session.Must(session.NewSession())
    cfg := aws.NewConfig().WithRegion(region)
    services, err := fetchEnabledServicesInRegion(&region, sess, cfg)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    activeServiceCount := len(services)

    fmt.Printf("Number of active services in region %s: %d\n", region, activeServiceCount)

    for service := range services {
        fmt.Println("Service:", service)
    }

    // Create an EC2 client
    ec2Svc := ec2.New(sess, cfg)

    // Use DescribeInstances to get information about running EC2 instances
    ec2Input := &ec2.DescribeInstancesInput{
        Filters: []*ec2.Filter{
            {
                Name:   aws.String("instance-state-name"),
                Values: []*string{aws.String("running")},
            },
        },
    }

    // ... your existing code for EC2 instance retrieval ...

    // Create a CloudWatch client
    cloudWatchSvc := cloudwatch.New(sess, cfg)

    // Specify the metric you want to retrieve (e.g., CPU utilization for EC2 instances)
    metricName := "CPUUtilization"
    namespace := "AWS/EC2"
    dimensions := []*cloudwatch.Dimension{
        {
            Name:  aws.String("InstanceId"),
            Value: aws.String("your-instance-id"), // Replace with your EC2 instance ID
        },
    }

    // Create a metric statistic query
    metricInput := &cloudwatch.GetMetricDataInput{
        MetricDataQueries: []*cloudwatch.MetricDataQuery{
            {
                Id: aws.String("m1"),
                MetricStat: &cloudwatch.MetricStat{
                    Metric: &cloudwatch.Metric{
                        Namespace:  aws.String(namespace),
                        MetricName: aws.String(metricName),
                        Dimensions: dimensions,
                    },
                    Period: aws.Int64(300), // Adjust the period as needed
                    Stat:   aws.String("Average"),
                },
                ReturnData: aws.Bool(true),
            },
        },
        StartTime: aws.Time(time.Now().Add(-1 * time.Hour)), // Adjust the start time as needed
        EndTime:   aws.Time(time.Now()),                    // Adjust the end time as needed
    }

    // Query the metric data
    metricData, err := cloudWatchSvc.GetMetricData(metricInput)
    if err != nil {
        fmt.Println("Error querying metric data:", err)
        return
    }

    // Process and display the metric data
    for _, result := range metricData.MetricDataResults {
        fmt.Printf("Metric Name: %s\n", *result.Label)
        for i, timestamp := range result.Timestamps {
            value := result.Values[i]
            fmt.Printf("Timestamp: %s, Value: %f\n", timestamp.Format(time.RFC3339), *value)
        }
    }

    // ... add code for billing information retrieval here ...
}

