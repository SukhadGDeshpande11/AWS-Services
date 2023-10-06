package main

import (
    "fmt"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/ssm"
    "strings"
    "unicode"
)

const ParamPath = "/aws/service/global-infrastructure/regions/%s/services"

func fetchEnabledServicesInRegion(region *string, sess *session.Session, cfg *aws.Config) (map[string]bool, error) {
    service := make(map[string]bool)
    svc := ssm.New(sess, cfg)
    var NextToken *string
next:
    req, resp := svc.GetParametersByPathRequest(&ssm.GetParametersByPathInput{
        Path:     aws.String(fmt.Sprintf(ParamPath, *region)),
        NextToken: NextToken,
    })

    err := req.Send()
    if err != nil {
        return nil, err
    }

    NextToken = resp.NextToken
    if len(resp.Parameters) > 0 {
        // fetch the service name, process
        for _, p := range resp.Parameters {
            if p.Value != nil {
                srv := formatServiceName(*p.Value)
                service[srv] = true
            }
        }
    }

    // AWS API call sometimes behave erratically, returning empty pages with NextToken as non null
    if NextToken != nil {
        goto next
    }
    return service, nil
}

func formatServiceName(s string) string {
    name := strings.TrimSpace(s)

    // Replace all Non-letter/number values with space
    // AWS services are not named consistently
    runes := []rune(name)
    for i := 0; i < len(runes); i++ {
        if r := runes[i]; !(unicode.IsNumber(r) || unicode.IsLetter(r)) {
            runes[i] = ' '
        }
    }
    name = string(runes)
    // Title case name so it's readable as a symbol.
    name = strings.ToLower(name)
    // Strip out spaces.
    name = strings.Replace(name, " ", "", -1)
    return name
}

func main() {
    // You should implement the main function as described in the previous responses.
    region := "ap-south-1" // Replace with your desired AWS region
    sess := session.Must(session.NewSession())
    cfg := aws.NewConfig().WithRegion(region)
    services, err := fetchEnabledServicesInRegion(&region, sess, cfg)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    for service := range services {
        fmt.Println("Service:", service)
    }
}

