package main

import (
    "fmt"
    "os"
    "time"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/costexplorer"
)

func main() {
    // Create an AWS session
    sess := session.Must(session.NewSession(&aws.Config{
        Region: aws.String("ap-south-1"),
    }))

    // Create a Cost Explorer client
    svc := costexplorer.New(sess)

    // Set the time range for which you want to fetch cost and usage data
    startDate := time.Now().AddDate(0, 0, -7) // 7 days ago
    endDate := time.Now()

    // Define your query input without grouping by service
    params := &costexplorer.GetCostAndUsageInput{
        TimePeriod: &costexplorer.DateInterval{
            Start: aws.String(startDate.Format("2023-09-02")),
            End:   aws.String(endDate.Format("2023-09-09")),
        },
        Granularity: aws.String("DAILY"), // You can change the granularity as needed (DAILY, MONTHLY, etc.)
        Metrics:     []*string{aws.String("UnblendedCost")}, // You can add more metrics as needed
    }

    // Fetch the cost and usage data
    result, err := svc.GetCostAndUsage(params)
    if err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }
    fmt.Println("cost and usage of the data",result)
}

