package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type AwsClient struct {
	ec2    *ec2.EC2
	output *string
}

type RawData struct {
	InstanceType string
	Role         string
	Environment  string
}

type Data struct {
	Environment   string
	InstanceType  string
	InstanceCount int
}

type CSVdata struct {
	Role string
	Data []Data
}

func NewAwsClient(region, output *string) *AwsClient {
	c := &AwsClient{
		ec2:    ec2.New(session.New(), &aws.Config{Region: region}),
		output: output,
	}
	return c
}

func (a *AwsClient) DescribeInstances() ([]RawData, error) {
	var rawDataSet []RawData
	rawData := RawData{}
	var nextToken *string

	for {
		input := &ec2.DescribeInstancesInput{
			NextToken: nextToken,
		}
		result, err := a.ec2.DescribeInstances(input)
		if err != nil {
			msg := fmt.Sprintf("Could not get instance details : %s", err.Error())
			describeInstancesErr := errors.New(msg)
			return nil, describeInstancesErr
		}
		nextToken = result.NextToken
		for _, reservation := range result.Reservations {
			for _, instance := range reservation.Instances {
				rawData.InstanceType = aws.StringValue(instance.InstanceType)
				for _, tag := range instance.Tags {
					if aws.StringValue(tag.Key) == "Environment" {
						rawData.Environment = aws.StringValue(tag.Value)
					} else if aws.StringValue(tag.Key) == "Role" {
						rawData.Role = aws.StringValue(tag.Value)
					}
				}
				rawDataSet = append(rawDataSet, rawData)
			}
		}
		if nextToken == nil {
			break
		}
	}
	return rawDataSet, nil
}

func GenerateCSVData(d []RawData, roles, environments []string) []CSVdata {
	var csvDataList []CSVdata

	for _, role := range roles {
		csvData := CSVdata{}
		tempCSVdata := CSVdata{}
		tempCSVdata.Role = role
		csvData.Role = role
		for _, data := range d {
			var tempData Data
			if data.Role == role {
				tempData.Environment = data.Environment
				tempData.InstanceType = data.InstanceType
				tempCSVdata.Data = append(tempCSVdata.Data, tempData)
			}
		}
		//fmt.Println(role, tempCSVdata.Data)

		for _, env := range environments {
			var instanceTypesForThisRoleInThisEnvironment []string
			var finalData Data
			for _, dt := range tempCSVdata.Data {
				if env == dt.Environment {
					instanceTypesForThisRoleInThisEnvironment = append(instanceTypesForThisRoleInThisEnvironment, dt.InstanceType)
					finalData.Environment = dt.Environment
					uniqueinstanceTypes := GenerateUniqueStringSlice(instanceTypesForThisRoleInThisEnvironment)

					for _, instanceType := range uniqueinstanceTypes {
						var t []string
						finalData.InstanceType = instanceType
						for _, i := range instanceTypesForThisRoleInThisEnvironment {
							if i == instanceType {
								t = append(t, i)
							}
						}
						finalData.InstanceCount = len(t)
					}
				}
			}
			if finalData.InstanceCount > 0 {
				csvData.Data = append(csvData.Data, finalData)
			}
		}
		csvDataList = append(csvDataList, csvData)
	}
	return csvDataList
}

func GenerateUniqueStringSlice(s []string) []string {
	u := make([]string, 0, len(s))
	m := make(map[string]bool)
	for _, val := range s {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}
	return u
}

func GetRoles(d []RawData) []string {
	roles := make([]string, 0, len(d))
	for _, rawData := range d {
		roles = append(roles, rawData.Role)
	}
	uniqueRoles := GenerateUniqueStringSlice(roles)
	return uniqueRoles
}

func GetEnvironments(d []RawData) []string {
	environments := make([]string, 0, len(d))
	for _, rawData := range d {
		environments = append(environments, rawData.Environment)
	}
	uniqueEnvironments := GenerateUniqueStringSlice(environments)
	return uniqueEnvironments
}

func GenerateCSV(d []CSVdata, output *string) {
	var records [][]string
	records = append(records, []string{"Role", "Environment", "InstanceType", "InstanceCount"})
	for _, data := range d {

		role := data.Role
		for _, dat := range data.Data {
			records = append(records, []string{role, dat.Environment, dat.InstanceType, strconv.Itoa(dat.InstanceCount)})
		}
	}
	if _, err := os.Stat(*output); os.IsNotExist(err) {
		os.MkdirAll(*output, os.ModePerm)

	}
	t := time.Now()
	timestamp := fmt.Sprintf("%02d-%02d-%d_%02d-%02d-%02d", t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute(), t.Second())
	outFilename := fmt.Sprintf("ec2instance_%s.csv", timestamp)
	outfile := filepath.Join(*output, outFilename)
	file, err := os.Create(outfile)

	if err != nil {
		log.Fatal("Could not create output csv", err)
	}
	defer file.Close()

	w := csv.NewWriter(file)

	w.WriteAll(records)
}

func (a *AwsClient) Handler() {
	rawDataSet, err := a.DescribeInstances()

	if err != nil {
		log.Fatal(err)
	}

	roles := GetRoles(rawDataSet)
	environments := GetEnvironments(rawDataSet)

	csvDataList := GenerateCSVData(rawDataSet, roles, environments)
	GenerateCSV(csvDataList, a.output)
}

func main() {
	awsRegionPtr := flag.String("region", "us-east-1", "AWS region")

	outputPtr := flag.String("output", "/tmp", "Output directory location")

	flag.Parse()

	awsClient := NewAwsClient(awsRegionPtr, outputPtr)

	awsClient.Handler()
}
