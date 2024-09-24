package biteship

type Request struct {
	OriginLatitude       float64 `json:"origin_latitude"`
	OriginLongitude      float64 `json:"origin_longitude"`
	DestinationLatitude  float64 `json:"destination_latitude"`
	DestinationLongitude float64 `json:"destination_longitude"`
	Couriers             string  `json:"couriers"`
	Items                []Items `json:"items"`
}

type Items struct {
	Name     string `json:"name"`
	Value    int64  `json:"value"`
	Weight   int32  `json:"weight"`
	Quantity int32  `json:"quantity"`
}

type ResponseShippingForProductDetail struct {
	CostRange        string `json:"cost_range"`
	EstimatedArrived string `json:"estimated_arrived"`
}

type Response struct {
	Success     bool             `json:"success"`
	Object      string           `json:"object"`
	Message     string           `json:"message"`
	Code        int              `json:"code"`
	Origin      Location         `json:"origin"`
	Stops       []interface{}    `json:"stops"`
	Destination Location         `json:"destination"`
	Pricing     []CourierPricing `json:"pricing"`
}

type Location struct {
	LocationID                       *string `json:"location_id"`
	Latitude                         float64 `json:"latitude"`
	Longitude                        float64 `json:"longitude"`
	PostalCode                       int     `json:"postal_code"`
	CountryName                      string  `json:"country_name"`
	CountryCode                      string  `json:"country_code"`
	AdministrativeDivisionLevel1Name string  `json:"administrative_division_level_1_name"`
	AdministrativeDivisionLevel1Type string  `json:"administrative_division_level_1_type"`
	AdministrativeDivisionLevel2Name string  `json:"administrative_division_level_2_name"`
	AdministrativeDivisionLevel2Type string  `json:"administrative_division_level_2_type"`
	AdministrativeDivisionLevel3Name string  `json:"administrative_division_level_3_name"`
	AdministrativeDivisionLevel3Type string  `json:"administrative_division_level_3_type"`
	AdministrativeDivisionLevel4Name string  `json:"administrative_division_level_4_name"`
	AdministrativeDivisionLevel4Type string  `json:"administrative_division_level_4_type"`
	Address                          string  `json:"address"`
}

type CourierPricing struct {
	AvailableCollectionMethod    []string `json:"available_collection_method"`
	AvailableForCashOnDelivery   bool     `json:"available_for_cash_on_delivery"`
	AvailableForProofOfDelivery  bool     `json:"available_for_proof_of_delivery"`
	AvailableForInstantWaybillID bool     `json:"available_for_instant_waybill_id"`
	AvailableForInsurance        bool     `json:"available_for_insurance"`
	Company                      string   `json:"company"`
	CourierName                  string   `json:"courier_name"`
	CourierCode                  string   `json:"courier_code"`
	CourierServiceName           string   `json:"courier_service_name"`
	CourierServiceCode           string   `json:"courier_service_code"`
	Description                  string   `json:"description"`
	Duration                     string   `json:"duration"`
	ShipmentDurationRange        string   `json:"shipment_duration_range"`
	ShipmentDurationUnit         string   `json:"shipment_duration_unit"`
	ServiceType                  string   `json:"service_type"`
	ShippingType                 string   `json:"shipping_type"`
	Price                        int      `json:"price"`
	Type                         string   `json:"type"`
}

const (
	LATITUDE  = -7.952240941097751
	LONGITUDE = 112.6126827676408
	COURIERS  = "gojek,grab,jne,sicepat,jnt,anteraja"
)
