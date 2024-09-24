package biteship

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func ShippingForProductDetails(latitude, longitude float64, items []Items) (ResponseShippingForProductDetail, error) {
	res := ResponseShippingForProductDetail{}

	url := "https://api.biteship.com/v1/rates/couriers"
	payload := Request{
		OriginLatitude:       LATITUDE,
		OriginLongitude:      LONGITUDE,
		Couriers:             COURIERS,
		DestinationLatitude:  latitude,
		DestinationLongitude: longitude,
		Items:                items,
	}

	p, err := json.Marshal(payload)
	if err != nil {
		return res, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(p))
	if err != nil {
		return res, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", os.Getenv("BITESHIP_API_KEY"))

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}

	var respStruct Response
	err = json.Unmarshal(body, &respStruct)
	if err != nil {
		return res, err
	}

	priceRange, timeRange := getPricingAndDurationRange(respStruct.Pricing)
	res.CostRange = priceRange
	res.EstimatedArrived = timeRange

	return res, nil
}

func getPricingAndDurationRange(pricings []CourierPricing) (string, string) {
	if len(pricings) == 0 {
		return "N/A", "N/A"
	}

	minPrice := pricings[0].Price
	maxPrice := pricings[0].Price
	minDays := 999999
	maxDays := 0

	for _, pricing := range pricings {
		if pricing.Price < minPrice {
			minPrice = pricing.Price
		}
		if pricing.Price > maxPrice {
			maxPrice = pricing.Price
		}

		shipmentRange := strings.Split(pricing.ShipmentDurationRange, " - ")
		if len(shipmentRange) == 2 {
			minShipmentDays, _ := strconv.Atoi(shipmentRange[0])
			maxShipmentDays, _ := strconv.Atoi(shipmentRange[1])

			if minShipmentDays < minDays {
				minDays = minShipmentDays
			}
			if maxShipmentDays > maxDays {
				maxDays = maxShipmentDays
			}
		} else if len(shipmentRange) == 1 {
			shipmentDays, _ := strconv.Atoi(shipmentRange[0])
			if shipmentDays < minDays {
				minDays = shipmentDays
			}
			if shipmentDays > maxDays {
				maxDays = shipmentDays
			}
		}
	}

	priceRange := fmt.Sprintf("Rp%d - Rp%d", minPrice, maxPrice)
	timeRange := fmt.Sprintf("%d - %d days", minDays, maxDays)

	return priceRange, timeRange
}
