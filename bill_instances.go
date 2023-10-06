package main

import (
    "fmt"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/ssm"
    "github.com/aws/aws-sdk-go/service/ec2"
    "github.com/aws/aws-sdk-go/service/cloudwatch"
    "strings"
    "unicode"
    "time"
)

const ParamPath = "/aws/service/global-infrastructure/regions/%s/services"

func fetchEnabledServicesInRegion(region *string, sess *session.Session, cfg *aws.Config) (map[string]bool, error) {
    // Existing code for fetching active services
    // ...

    return service, nil
}

func fetchEC2InstancesInRegion(sess *session.Session, cfg *aws.Config) ([]*ec2.Instance, error) {
    ec2Svc := ec2.New(sess, cfg)

    ec2Input := &ec2.DescribeInstancesInput{
        Filters: []*ec2.Filter{
            {
                Name:   aws.String("instance-state-name"),
                Values: []*string{aws.String("running")},
            },
        },
    }

    ec2Result, err := ec2Svc.DescribeInstances(ec2Input)
    if err != nil {
        return nil, err
    }

    instances := []*ec2.Instance{}
    for _, reservation := range ec2Result.Reservations {
        for _, instance := range reservation.Instances {
            instances = append(instances, instance)
        }
    }
    return instances, nil
}

func fetchEC2CPUUtilization(instanceID string, cloudWatchSvc *cloudwatch.CloudWatch, startTime, endTime time.Time) (float64, error) {
    // Define the metric parameters
    metricName := "CPUUtilization"
    namespace := "AWS/EC2"
    dimensions := []*cloudwatch.Dimension{
        {
            Name:  aws.String("InstanceId"),
            Value: aws.String(instanceID),
        },
    }

    // Define the metric query
    metricDataQuery := &cloudwatch.MetricDataQuery{
        Id:         aws.String("m1"),
        MetricStat: &cloudwatch.MetricStat{},
    }
    metricDataQuery.MetricStat.Metric = &cloudwatch.Metric{
        MetricName: aws.String(metricName),
        Namespace:  aws.String(namespace),
        Dimensions: dimensions,
    }
    metricDataQuery.MetricStat.Period = aws.Int64(300) // Adjust as needed
    metricDataQuery.MetricStat.Stat = aws.String("Average")
    metricDataQuery.MetricStat.Unit = aws.String("Percent")

    // Fetch the metric data
    metricDataInput := &cloudwatch.GetMetricDataInput{
        MetricDataQueries: []*cloudwatch.MetricDataQuery{metricDataQuery},
        StartTime:         aws.Time(startTime),
        EndTime:           aws.Time(endTime),
    }

    metricDataOutput, err := cloudWatchSvc.GetMetricData(metricDataInput)
    if err != nil {
        return 0.0, err
    }

    // Extract and return the CPU utilization value
    if len(metricDataOutput.MetricDataResults) > 0 && len(metricDataOutput.MetricDataResults[0].Values) > 0 {
        return *metricDataOutput.MetricDataResults[0].Values[0], nil
    }

    return 0.0, fmt.Errorf("No CPU utilization data available")
}

func main() {
    region := "ap-south-1"
    sess := session.Must(session.NewSession())
    cfg := aws.NewConfig().WithRegion(region)

    // Fetch active services
    services, err := fetchEnabledServicesInRegion(&region, sess, cfg)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    // Fetch active EC2 instances
    instances, err := fetchEC2InstancesInRegion(sess, cfg)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    // Create a CloudWatch client
    cloudWatchSvc := cloudwatch.New(sess, cfg)

    // Set the time range for the metric data (e.g., last hour)
    endTime := time.Now()
    startTime := endTime.Add(-time.Hour)

    // Fetch and print CPU utilization for each active EC2 instance
    for _, instance := range instances {
        instanceID := *instance.InstanceId
        cpuUtilization, err := fetchEC2CPUUtilization(instanceID, cloudWatchSvc, startTime, endTime)
        if err != nil {
            fmt.Printf("Error fetching CPU utilization for instance %s: %v\n", instanceID, err)
        } else {
            fmt.Printf("Instance ID: %s, CPU Utilization: %.2f%%\n", instanceID, cpuUtilization)
        }
    }
}

