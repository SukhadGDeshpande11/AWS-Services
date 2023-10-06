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
		for _, p := range resp.Parameters {
			if p.Value != nil {
				srv := formatServiceName(*p.Value)
				service[srv] = true
			}
		}
	}

	if NextToken != nil {
		goto next
	}
	return service, nil
}

func formatServiceName(s string) string {
	name := strings.TrimSpace(s)
	runes := []rune(name)
	for i := 0; i < len(runes); i++ {
		if r := runes[i]; !(unicode.IsNumber(r) || unicode.IsLetter(r)) {
			runes[i] = ' '
		}
	}
	name = string(runes)
	name = strings.ToLower(name)
	name = strings.Replace(name, " ", "", -1)
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
}

