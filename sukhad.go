package main

import (
    "fmt"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/ec2"
    "github.com/aws/aws-sdk-go/service/cloudwatch"
    "time"
)

const ParamPath = "/aws/service/global-infrastructure/regions/%s/services"

func fetchEnabledServicesInRegion(region *string, sess *session.Session, cfg *aws.Config) (map[string]bool, error){

     return nil, nil
}//

func fetchEC2InstancesInRegion(sess *session.Session, cfg *aws.Config) ([]*ec2.Instance, error) {
    ec2Svc := ec2.New(sess, cfg)

    ec2Input := &ec2.DescribeInstancesInput{    //ec2 .describe instancesinput filters out and lists only running instances
        Filters: []*ec2.Filter{
            {
                Name:   aws.String("instance-state-name"),
                Values: []*string{aws.String("running")},
            },
        },
    }

    ec2Result, err := ec2Svc.DescribeInstances(ec2Input)// over here ec2 service with the ip it matches the servies 
    if err != nil {
        return nil, err
    }

    instances := []*ec2.Instance{}
    for _, reservation := range ec2Result.Reservations {
        for _, instance := range reservation.Instances {
            instances = append(instances, instance)
        }
    }
    return instances, nil// if the api call is successfull it initialise an empty slice of ec.2 instances so basically a
    			// for loop has been used to iterate these instances 
}// over here it takes a session as an object that represents aws session cfg configures the session 

func fetchEC2CPUUtilization(instanceID string, cloudWatchSvc *cloudwatch.CloudWatch, startTime, endTime time.Time) (float64, error) {

    return 0.0, nil 
}

func fetchEC2BillingInfo(cloudWatchSvc *cloudwatch.CloudWatch, startTime, endTime time.Time) (float64, error) {
    metricName := "EstimatedCharges"//? what metric name should i give here 
    namespace := "AWS/Billing"//?? what namespace should i giver here 
    dimensions := []*cloudwatch.Dimension{
        {
            Name:  aws.String("Currency"),
            Value: aws.String("USD"),
        },// what details i should give so that it works properly
    }

    metricDataQuery := &cloudwatch.MetricDataQuery{ // creating the metric data query that will give me to retrive the inforation 
        Id:         aws.String("m1"),
        MetricStat: &cloudwatch.MetricStat{},
    }
    metricDataQuery.MetricStat.Metric = &cloudwatch.Metric{
        MetricName: aws.String(metricName),
        Namespace:  aws.String(namespace),
        Dimensions: dimensions,
    }
    metricDataQuery.MetricStat.Period = aws.Int64(86400) // setting metric data query for 24 hours converted into seconds
    metricDataQuery.MetricStat.Stat = aws.String("Maximum")

    metricDataInput := &cloudwatch.GetMetricDataInput{
        MetricDataQueries: []*cloudwatch.MetricDataQuery{metricDataQuery},
        StartTime:         aws.Time(startTime),
        EndTime:           aws.Time(endTime),
    }// need to give the start time and end time over here 

    metricDataOutput, err := cloudWatchSvc.GetMetricData(metricDataInput)
    if err != nil {
        return 0.0, err
    }
    // checking if there are Matric data results in the response if yes then return the op if not nill

    if len(metricDataOutput.MetricDataResults) > 0 && len(metricDataOutput.MetricDataResults[0].Values) > 0 {
        return *metricDataOutput.MetricDataResults[0].Values[0], nil
    }

    return 0.0, fmt.Errorf("No billing information available")
}

func main() {
    region := "ap-south-1"// giving aws region
    sess := session.Must(session.NewSession()) //creating new session and confirguing it 
    cfg := aws.NewConfig().WithRegion(region)
    instances, err := fetchEC2InstancesInRegion(sess, cfg) // calling ec2 instances and fetching the services 
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    cloudWatchSvc := cloudwatch.New(sess, cfg) // creating a cloudwatch here 
    endTime := time.Now()
    startTime := endTime.Add(-24 * time.Hour) // setting the start and the end time here 
    for _, instance := range instances { // iterationg with respect to ec2 instances 
        instanceID := *instance.InstanceId // fetching the instances id 
        cpuUtilization, err := fetchEC2CPUUtilization(instanceID, cloudWatchSvc, startTime, endTime)// computing cpu 
        if err != nil {
            fmt.Printf("Error fetching CPU utilization for instance %s: %v\n", instanceID, err)
        } else {
            fmt.Printf("Instance ID: %s, CPU Utilization: %.2f%%\n", instanceID, cpuUtilization)
        }
		billingInfo, err := fetchEC2BillingInfo(cloudWatchSvc, startTime, endTime)
        if err != nil {
            fmt.Printf("Error fetching billing info for instance %s: %v\n", instanceID, err)
        } else {
            fmt.Printf("Instance ID: %s, Billing Info (USD): %.2f\n", instanceID, billingInfo)
        }
    }
}

