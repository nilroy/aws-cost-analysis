package main

import (
	"bytes"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type AwsClient struct {
	ec2     *ec2.EC2
	s3      *s3.S3
	region  *string
	tempDir *string
	s3Dir   *string
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

func NewAwsClient(region, tempDir, s3Dir *string) *AwsClient {
	c := &AwsClient{
		ec2:     ec2.New(session.New(), &aws.Config{Region: region}),
		s3:      s3.New(session.New(), &aws.Config{Region: region}),
		tempDir: tempDir,
		s3Dir:   s3Dir,
		region:  region,
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

				if aws.StringValue(instance.State.Name) != "running" {
					continue
				}

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
		for _, rawData := range d {
			var tempData Data
			if rawData.Role == role {
				tempData.Environment = rawData.Environment
				tempData.InstanceType = rawData.InstanceType
				tempCSVdata.Data = append(tempCSVdata.Data, tempData)
			}
		}

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

func GenerateCSV(d []CSVdata, tempDir *string, ec2PriceList map[string]EC2InstancePrice) (tempFilePath string) {
	var records [][]string
	records = append(records, []string{
		"Role",
		"Environment",
		"InstanceType",
		"InstanceCount",
		"OnDemand",
		"YrTerm1ConvertibleAllUpfront",
		"YrTerm1ConvertiblePartialUpfront",
		"YrTerm1ConvertibleNoUpfront",
		"YrTerm1StandardAllUpfront",
		"YrTerm1StandardPartialUpfront",
		"YrTerm1StandardNoUpfront",
		"YrTerm3ConvertibleAllUpfront",
		"YrTerm3ConvertiblePartialUpfront",
		"YrTerm3ConvertibleNoUpfront",
		"YrTerm3StandardAllUpfront",
		"YrTerm3StandardPartialUpfront",
		"YrTerm3StandardNoUpfront",
	})
	for _, data := range d {

		role := data.Role
		for _, dat := range data.Data {
			instanceType := dat.InstanceType
			instanceCount := strconv.Itoa(dat.InstanceCount)

			// All prices calculated on yearly basis

			totalUsageHoursPerYear := float64(24*365) * float64(dat.InstanceCount)

			OnDemand := fmt.Sprintf("%.2f", *ec2PriceList[instanceType].OnDemand*totalUsageHoursPerYear)

			YrTerm1ConvertibleAllUpfront := fmt.Sprintf("%.2f", *ec2PriceList[instanceType].YrTerm1ConvertibleAllUpfront*totalUsageHoursPerYear)
			YrTerm1ConvertiblePartialUpfront := fmt.Sprintf("%.2f", *ec2PriceList[instanceType].YrTerm1ConvertiblePartialUpfront*totalUsageHoursPerYear)
			YrTerm1ConvertibleNoUpfront := fmt.Sprintf("%.2f", *ec2PriceList[instanceType].YrTerm1ConvertibleNoUpfront*totalUsageHoursPerYear)

			YrTerm1StandardAllUpfront := fmt.Sprintf("%.2f", *ec2PriceList[instanceType].YrTerm1StandardAllUpfront*totalUsageHoursPerYear)
			YrTerm1StandardPartialUpfront := fmt.Sprintf("%.2f", *ec2PriceList[instanceType].YrTerm1StandardPartialUpfront*totalUsageHoursPerYear)
			YrTerm1StandardNoUpfront := fmt.Sprintf("%.2f", *ec2PriceList[instanceType].YrTerm1StandardNoUpfront*totalUsageHoursPerYear)

			YrTerm3ConvertibleAllUpfront := fmt.Sprintf("%.2f", *ec2PriceList[instanceType].YrTerm3ConvertibleAllUpfront*totalUsageHoursPerYear)
			YrTerm3ConvertiblePartialUpfront := fmt.Sprintf("%.2f", *ec2PriceList[instanceType].YrTerm3ConvertiblePartialUpfront*totalUsageHoursPerYear)
			YrTerm3ConvertibleNoUpfront := fmt.Sprintf("%.2f", *ec2PriceList[instanceType].YrTerm3ConvertibleNoUpfront*totalUsageHoursPerYear)

			YrTerm3StandardAllUpfront := fmt.Sprintf("%.2f", *ec2PriceList[instanceType].YrTerm3StandardAllUpfront*totalUsageHoursPerYear)
			YrTerm3StandardPartialUpfront := fmt.Sprintf("%.2f", *ec2PriceList[instanceType].YrTerm3StandardPartialUpfront*totalUsageHoursPerYear)
			YrTerm3StandardNoUpfront := fmt.Sprintf("%.2f", *ec2PriceList[instanceType].YrTerm3StandardNoUpfront*totalUsageHoursPerYear)

			records = append(records, []string{
				role,
				dat.Environment,
				instanceType,
				instanceCount,
				OnDemand,
				YrTerm1ConvertibleAllUpfront,
				YrTerm1ConvertiblePartialUpfront,
				YrTerm1ConvertibleNoUpfront,
				YrTerm1StandardAllUpfront,
				YrTerm1StandardPartialUpfront,
				YrTerm1StandardNoUpfront,
				YrTerm3ConvertibleAllUpfront,
				YrTerm3ConvertiblePartialUpfront,
				YrTerm3ConvertibleNoUpfront,
				YrTerm3StandardAllUpfront,
				YrTerm3StandardPartialUpfront,
				YrTerm3StandardNoUpfront,
			})
		}
	}
	if _, err := os.Stat(*tempDir); os.IsNotExist(err) {
		os.MkdirAll(*tempDir, os.ModePerm)

	}
	t := time.Now()
	timestamp := fmt.Sprintf("%02d-%02d-%d_%02d-%02d-%02d", t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute(), t.Second())
	outFilename := fmt.Sprintf("ec2instance_%s.csv", timestamp)
	outfile := filepath.Join(*tempDir, outFilename)
	file, err := os.Create(outfile)

	if err != nil {
		log.Fatal("Could not create temp csv", err)
	}
	defer file.Close()

	w := csv.NewWriter(file)

	w.WriteAll(records)
	return outfile
}

func (a *AwsClient) UploadToS3(filePath string) error {
	file, err := os.Open(filePath)

	if err != nil {
		return err
	}
	defer file.Close()
	fileInfo, _ := file.Stat()
	size := fileInfo.Size()

	buffer := make([]byte, size)
	file.Read(buffer)
	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)
	path := "/ec2/nosto-ec2-instance-details.csv"

	input := &s3.PutObjectInput{
		Body:          fileBytes,
		Bucket:        a.s3Dir,
		Key:           aws.String(path),
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(fileType),
	}
	_, err = a.s3.PutObject(input)
	return err
}

func (a *AwsClient) Handler() {
	ec2PriceList := GeneratePriceList(*a.region)
	rawDataSet, err := a.DescribeInstances()

	if err != nil {
		log.Fatal(err)
	}

	roles := GetRoles(rawDataSet)
	environments := GetEnvironments(rawDataSet)

	csvDataList := GenerateCSVData(rawDataSet, roles, environments)

	tempFile := GenerateCSV(csvDataList, a.tempDir, ec2PriceList)
	err = a.UploadToS3(tempFile)

	if err != nil {
		log.Fatal(err)
	}

	err = os.RemoveAll(tempFile)

	if err != nil {
		log.Fatal(err)
	}

}

func usage() {
	fmt.Fprintf(os.Stderr, "usage : %s -region <region> -s3 <s3 bucket name>", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	awsRegionPtr := flag.String("region", "us-east-1", "AWS region")

	tempPtr := flag.String("temp", "/tmp/ec2pricing", "temp directory location")

	s3Ptr := flag.String("s3", "", "S3 upload location")

	flag.Parse()

	if *s3Ptr == "" {
		usage()
		os.Exit(1)
	}

	awsClient := NewAwsClient(awsRegionPtr, tempPtr, s3Ptr)

	awsClient.Handler()
}
