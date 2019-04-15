package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type ReservedInstancePrice struct {
	YrTerm1ConvertibleAllUpfront     *string `json:"yrTerm1Convertible.allUpfront"`
	YrTerm1ConvertiblePartialUpfront *string `json:"yrTerm1Convertible.partialUpfront"`
	YrTerm1ConvertibleNoUpfront      *string `json:"yrTerm1Convertible.noUpfront"`
	YrTerm1StandardAllUpfront        *string `json:"yrTerm1Standard.allUpfront"`
	YrTerm1StandardPartialUpfront    *string `json:"yrTerm1Standard.partialUpfront"`
	YrTerm1StandardNoUpfront         *string `json:"yrTerm1Standard.noUpfront"`
	YrTerm3ConvertibleAllUpfront     *string `json:"yrTerm3Convertible.allUpfront"`
	YrTerm3ConvertiblePartialUpfront *string `json:"yrTerm3Convertible.partialUpfront"`
	YrTerm3ConvertibleNoUpfront      *string `json:"yrTerm3Convertible.noUpfront"`
	YrTerm3StandardAllUpfront        *string `json:"yrTerm3Standard.allUpfront"`
	YrTerm3StandardPartialUpfront    *string `json:"yrTerm3Standard.partialUpfront"`
	YrTerm3StandardNoUpfront         *string `json:"yrTerm3Standard.noUpfront"`
}

type LinuxPriceDetail struct {
	OnDemand *string                `json:"ondemand"`
	Reserved *ReservedInstancePrice `json:"reserved"`
}

type PricingData struct {
	LinuxPrice *LinuxPriceDetail `json:"linux"`
}

type PricingDetails struct {
	InstanceType string                 `json:"instance_type"`
	Pricing      map[string]PricingData `json:"pricing"`
}

type EC2InstancePrice struct {
	OnDemand                         *float64
	YrTerm1ConvertibleAllUpfront     *float64
	YrTerm1ConvertiblePartialUpfront *float64
	YrTerm1ConvertibleNoUpfront      *float64
	YrTerm1StandardAllUpfront        *float64
	YrTerm1StandardPartialUpfront    *float64
	YrTerm1StandardNoUpfront         *float64
	YrTerm3ConvertibleAllUpfront     *float64
	YrTerm3ConvertiblePartialUpfront *float64
	YrTerm3ConvertibleNoUpfront      *float64
	YrTerm3StandardAllUpfront        *float64
	YrTerm3StandardPartialUpfront    *float64
	YrTerm3StandardNoUpfront         *float64
}

func GeneratePriceList(awsRegion string) map[string]EC2InstancePrice {
	url := "https://raw.githubusercontent.com/powdahound/ec2instances.info/master/www/instances.json"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var p []PricingDetails
	err = json.Unmarshal(body, &p)
	if err != nil {
		fmt.Println(err)
	}
	ec2InstancePriceList := map[string]EC2InstancePrice{}

	for _, priceDetail := range p {
		var ec2InstancePrice EC2InstancePrice

		if len(priceDetail.Pricing) != 0 {
			if priceDetail.Pricing[awsRegion].LinuxPrice.OnDemand != nil {
				onDemandPrice, _ := strconv.ParseFloat(*priceDetail.Pricing[awsRegion].LinuxPrice.OnDemand, 64)
				ec2InstancePrice.OnDemand = &onDemandPrice
			}

			if priceDetail.Pricing[awsRegion].LinuxPrice.Reserved != nil {
				if priceDetail.Pricing[awsRegion].LinuxPrice.Reserved.YrTerm1ConvertibleAllUpfront != nil {
					YrTerm1ConvertibleAllUpfront, _ := strconv.ParseFloat(*priceDetail.Pricing[awsRegion].LinuxPrice.Reserved.YrTerm1ConvertibleAllUpfront, 64)
					ec2InstancePrice.YrTerm1ConvertibleAllUpfront = &YrTerm1ConvertibleAllUpfront
				}

				if priceDetail.Pricing[awsRegion].LinuxPrice.Reserved.YrTerm1ConvertibleNoUpfront != nil {
					YrTerm1ConvertibleNoUpfront, _ := strconv.ParseFloat(*priceDetail.Pricing[awsRegion].LinuxPrice.Reserved.YrTerm1ConvertibleNoUpfront, 64)
					ec2InstancePrice.YrTerm1ConvertibleNoUpfront = &YrTerm1ConvertibleNoUpfront
				}

				if priceDetail.Pricing[awsRegion].LinuxPrice.Reserved.YrTerm1ConvertiblePartialUpfront != nil {
					YrTerm1ConvertiblePartialUpfront, _ := strconv.ParseFloat(*priceDetail.Pricing[awsRegion].LinuxPrice.Reserved.YrTerm1ConvertiblePartialUpfront, 64)
					ec2InstancePrice.YrTerm1ConvertiblePartialUpfront = &YrTerm1ConvertiblePartialUpfront
				}

				if priceDetail.Pricing[awsRegion].LinuxPrice.Reserved.YrTerm1StandardAllUpfront != nil {
					YrTerm1StandardAllUpfront, _ := strconv.ParseFloat(*priceDetail.Pricing[awsRegion].LinuxPrice.Reserved.YrTerm1StandardAllUpfront, 64)
					ec2InstancePrice.YrTerm1StandardAllUpfront = &YrTerm1StandardAllUpfront
				}

				if priceDetail.Pricing[awsRegion].LinuxPrice.Reserved.YrTerm1StandardNoUpfront != nil {
					YrTerm1StandardNoUpfront, _ := strconv.ParseFloat(*priceDetail.Pricing[awsRegion].LinuxPrice.Reserved.YrTerm1StandardNoUpfront, 64)
					ec2InstancePrice.YrTerm1StandardNoUpfront = &YrTerm1StandardNoUpfront
				}

				if priceDetail.Pricing[awsRegion].LinuxPrice.Reserved.YrTerm1StandardPartialUpfront != nil {
					YrTerm1StandardPartialUpfront, _ := strconv.ParseFloat(*priceDetail.Pricing[awsRegion].LinuxPrice.Reserved.YrTerm1StandardPartialUpfront, 64)
					ec2InstancePrice.YrTerm1StandardPartialUpfront = &YrTerm1StandardPartialUpfront
				}

				if priceDetail.Pricing[awsRegion].LinuxPrice.Reserved.YrTerm3ConvertibleAllUpfront != nil {
					YrTerm3ConvertibleAllUpfront, _ := strconv.ParseFloat(*priceDetail.Pricing[awsRegion].LinuxPrice.Reserved.YrTerm3ConvertibleAllUpfront, 64)
					ec2InstancePrice.YrTerm3ConvertibleAllUpfront = &YrTerm3ConvertibleAllUpfront
				}

				if priceDetail.Pricing[awsRegion].LinuxPrice.Reserved.YrTerm3ConvertibleNoUpfront != nil {
					YrTerm3ConvertibleNoUpfront, _ := strconv.ParseFloat(*priceDetail.Pricing[awsRegion].LinuxPrice.Reserved.YrTerm3ConvertibleNoUpfront, 64)
					ec2InstancePrice.YrTerm3ConvertibleNoUpfront = &YrTerm3ConvertibleNoUpfront
				}

				if priceDetail.Pricing[awsRegion].LinuxPrice.Reserved.YrTerm3ConvertiblePartialUpfront != nil {
					YrTerm3ConvertiblePartialUpfront, _ := strconv.ParseFloat(*priceDetail.Pricing[awsRegion].LinuxPrice.Reserved.YrTerm3ConvertiblePartialUpfront, 64)
					ec2InstancePrice.YrTerm3ConvertiblePartialUpfront = &YrTerm3ConvertiblePartialUpfront
				}

				if priceDetail.Pricing[awsRegion].LinuxPrice.Reserved.YrTerm3StandardAllUpfront != nil {
					YrTerm3StandardAllUpfront, _ := strconv.ParseFloat(*priceDetail.Pricing[awsRegion].LinuxPrice.Reserved.YrTerm3StandardAllUpfront, 64)
					ec2InstancePrice.YrTerm3StandardAllUpfront = &YrTerm3StandardAllUpfront
				}

				if priceDetail.Pricing[awsRegion].LinuxPrice.Reserved.YrTerm3StandardNoUpfront != nil {
					YrTerm3StandardNoUpfront, _ := strconv.ParseFloat(*priceDetail.Pricing[awsRegion].LinuxPrice.Reserved.YrTerm3StandardNoUpfront, 64)
					ec2InstancePrice.YrTerm3StandardNoUpfront = &YrTerm3StandardNoUpfront
				}

				if priceDetail.Pricing[awsRegion].LinuxPrice.Reserved.YrTerm3StandardPartialUpfront != nil {
					YrTerm3StandardPartialUpfront, _ := strconv.ParseFloat(*priceDetail.Pricing[awsRegion].LinuxPrice.Reserved.YrTerm3StandardPartialUpfront, 64)
					ec2InstancePrice.YrTerm3StandardPartialUpfront = &YrTerm3StandardPartialUpfront
				}

			}

			ec2InstancePriceList[priceDetail.InstanceType] = ec2InstancePrice

		}
	}
	return ec2InstancePriceList
}
