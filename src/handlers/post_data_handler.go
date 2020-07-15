package handlers

import (
	"github.com/ugorji/go/codec"
	"go-api/src/models"
	"go-api/src/repository"
	"encoding/json"
	"log"
)

func PostCustomerCompanies(result [][]string, r *repository.Repository) error {
	var row []string
	var customerCompanies []models.CustomerCompanies
	var customerCompany models.CustomerCompanies
	for i := 1 ; i<len(result) ; i++ {
		row = DecodeFileToString(result[i])

		b, err := json.Marshal(row)
		if err != nil {
			log.Fatal("\nERROR while marshaling the data ", err)
		}
		codec.NewDecoderBytes(b, new(codec.JsonHandle)).Decode(&customerCompany)
		customerCompanies = append(customerCompanies, customerCompany)
	}

	err := r.CreateCustomerCompanies(customerCompanies)
	return err
}

func PostOrderItemsData(result [][]string, r *repository.Repository) error {
	var row []string
	var orderItems []models.OrderItems
	var orderItem models.OrderItems
	for i := 1 ; i<len(result) ; i++ {
		row = DecodeFileToString(result[i])
		b, err := json.Marshal(row)
		if err != nil {
			log.Fatal("\nERROR while marshaling the data ", err)
		}
		codec.NewDecoderBytes(b, new(codec.JsonHandle)).Decode(&orderItem)
		orderItems = append(orderItems, orderItem)
	}

	err := r.CreateOrderItems(orderItems)
	return err
}

func PostDeliveriesData(result [][]string, r *repository.Repository) error {
	var row []string
	var deliveries []models.Deliveries
	var delivery models.Deliveries
	for i := 1 ; i<len(result) ; i++ {
		row = DecodeFileToString(result[i])

		b, err := json.Marshal(row)
		if err != nil {
			log.Fatal("\nERROR while marshaling the data ", err)
		}
		codec.NewDecoderBytes(b, new(codec.JsonHandle)).Decode(&delivery)
		deliveries = append(deliveries, delivery)
	}

	err := r.CreateDeliveries(deliveries)
	return err
}

func PostCustomerData(result [][]string, r *repository.Repository) error {
	var row []string
	var customers []models.Customers
	var customer models.Customers
	var creditCards []string
	
	for i := 1 ; i<len(result) ; i++ {
		row = DecodeFileToString(result[i])
		b, err := json.Marshal(row)
		if err != nil {
			log.Fatal("\nERROR while marshaling the data ", err)
		}
		
        creditCards = append(creditCards, result[i][len(result[i]) - 1])
		codec.NewDecoderBytes(b, new(codec.JsonHandle)).Decode(&customer)
		customers = append(customers, customer)
	}
	err := r.CreateCustomers(customers, creditCards)
	return err
}

func PostOrderData(result [][]string, r *repository.Repository) error {
	var row []string
	var orders []models.Orders
	var order models.Orders
	for i := 1 ; i<len(result) ; i++ {
		row = DecodeFileToString(result[i])

		b, err := json.Marshal(row)
		if err != nil {
			log.Fatal("\nERROR while marshaling the data ", err)
		}
		codec.NewDecoderBytes(b, new(codec.JsonHandle)).Decode(&order)
		orders = append(orders, order)
	}

	err := r.CreateOrders(orders)
	return err
}
