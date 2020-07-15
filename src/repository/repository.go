package repository

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"go-api/src/models"
	"log"
	"net/http"
	"go-api/src/response"
	"strings"
)

type Repository struct {
	DB         *sql.DB
	DBUri      string
	DBName     string
	DBUsername string
	DBPassword string
	Table      string
}

type Dberror struct {
	Code int
	Err  string
}

func (e Dberror) Error() string {
	return e.Err
}

func (r *Repository) Init() error {
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", r.DBUsername, r.DBPassword, r.DBUri, r.DBName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}
	log.Print("db client initialized")
	r.DB = db
	return nil
}

func (r Repository) CloseConn() error {
	return r.DB.Close()
}

func (r Repository) ListOrders() ([]response.OrdersResponse, error) {
	var query = "select o.order_name, cc.company_name, c.name, o.created_at, d.delivered_quantity*oi.price_per_unit as delivered_amount, oi.quantity*oi.price_per_unit as total_amount " +
		"from order_items oi, orders o, deliveries d, customer_companies cc, customers c " +
		"where oi.order_id = o.id and o.customer_id = c.user_id and c.company_id = cc.company_id and d.order_item = oi.order_id;"
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, Dberror{
			Code: http.StatusInternalServerError,
			Err:  err.Error(),
		}
	}

	orders := make([]response.OrdersResponse, 0)

	for rows.Next() {
		order := response.OrdersResponse{}
		err := rows.Scan(&order.OrderName, &order.CustomerCompany, &order.CustomerName, &order.OrderDate, &order.DeliveredAmount, &order.TotalAmount)
		if err != nil {
			return nil, Dberror{
				Code: http.StatusInternalServerError,
				Err:  fmt.Sprintf("error in conversion from db row: %s", err.Error()),
			}
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (r Repository) GetOrder(orderName string) ([]response.OrdersResponse, error) {
	dates:= strings.Split(orderName,",")
	if len(dates) != 2 {
		return nil, nil
	}
	var query = "select o.order_name, cc.company_name, c.name, o.created_at, d.delivered_quantity*oi.price_per_unit as delivered_amount, oi.quantity*oi.price_per_unit as total_amount " +
		"from order_items oi, orders o, deliveries d, customer_companies cc, customers c " +
		"where oi.order_id = o.id and o.customer_id = c.user_id and c.company_id = cc.company_id and d.order_item = oi.order_id " +
		"and o.created_at >= '" + dates[0] +"' and o.created_at <= '"+ dates[1] +"';"
	fmt.Println(query)
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, Dberror{
			Code: http.StatusInternalServerError,
			Err:  err.Error(),
		}
	}
	orders := make([]response.OrdersResponse, 0)
	for rows.Next() {
		order := response.OrdersResponse{}
		err := rows.Scan(&order.OrderName, &order.CustomerCompany, &order.CustomerName, &order.OrderDate, &order.DeliveredAmount, &order.TotalAmount)
		if err != nil {
			return nil, Dberror{
				Code: http.StatusInternalServerError,
				Err:  fmt.Sprintf("error in conversion from db row: %s", err.Error()),
			}
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (r Repository) CreateOrders(orders []models.Orders) error {
	var order models.Orders
	insertQry := "insert into orders (id, created_at, order_name, customer_id) values ($1, $2, $3, $4);"
	if len(orders) != 0 {
		for _,order = range orders {
			_, err := r.DB.Exec(insertQry, order.ID, order.CREATED_AT, order.ORDER_NAME, order.CUSTOMER_ID)
			if err != nil {
				return Dberror{
					Code: http.StatusInternalServerError,
					Err:  fmt.Sprintf("error in inserting the data in postgres: %s", err.Error()),
				}
			}
		}
	}
	return nil
}

func (r *Repository) CreateCustomers(customers []models.Customers, creditCards []string) error {
	var customer models.Customers
    i := 0       
	insertQry := "insert into customers (user_id, login, password, name, company_id, credit_cards) values ($1, $2, $3, $4, $5, $6);"
	if len(customers) != 0 {
		for _,customer = range customers {
			creditCards := Parse(creditCards[i])
			_, err := r.DB.Exec(insertQry, customer.USER_ID, customer.LOGIN, customer.PASSWORD, customer.NAME, customer.COMPANY_ID, pq.Array(creditCards))
			if err != nil {
				return Dberror{
					Code: http.StatusInternalServerError,
					Err:  fmt.Sprintf("error in inserting the data in postgres: %s", err.Error()),
				}
			}
			i++
		}
	}
	return nil
}

 func (r *Repository) CreateOrderItems (orderItems []models.OrderItems) error {
	var orderItem models.OrderItems
	
	insertQry := "insert into order_items (id, order_id, price_per_unit, quantity, product) values ($1, $2, $3, $4, $5);"
	if len(orderItems) != 0 {
		for _,orderItem = range orderItems {
			_, err := r.DB.Exec(insertQry, orderItem.ID, orderItem.ORDER_ID, orderItem.PRICE_PER_UNIT, orderItem.QUANTITY, orderItem.PRODUCT)
			if err != nil {
				return Dberror{
					Code: http.StatusInternalServerError,
					Err:  fmt.Sprintf("error in inserting the data in postgres: %s", err.Error()),
				}
			}
		}
	}
	return nil
}

func (r *Repository) CreateDeliveries (deliveries []models.Deliveries) error {
	var delivery models.Deliveries
	
	insertQry := "insert into deliveries (id, order_item, delivered_quantity) values ($1, $2, $3);"
	if len(deliveries) != 0 {
		for _,delivery = range deliveries {
			_, err := r.DB.Exec(insertQry, delivery.ID, delivery.ORDER_ITEM, delivery.DELIVERED_QUANTITY)
			if err != nil {
				return Dberror{
					Code: http.StatusInternalServerError,
					Err:  fmt.Sprintf("error in inserting the data in postgres: %s", err.Error()),
				}
			}
		}
	}
	return nil
}

func (r *Repository) CreateCustomerCompanies (customerCompanies []models.CustomerCompanies) error {
	var customerCompany models.CustomerCompanies
	
	insertQry := "insert into customer_companies (company_id, company_name) values ($1, $2);"
	if len(customerCompanies) != 0 {
		for _,customerCompany = range customerCompanies {
			_, err := r.DB.Exec(insertQry, customerCompany.COMPANY_ID, customerCompany.COMPANY_NAME)
			if err != nil {
				return Dberror{
					Code: http.StatusInternalServerError,
					Err:  fmt.Sprintf("error in inserting the data in postgres: %s", err.Error()),
				}
			}
		}
	}
	return nil
}

func Parse(str string) []string {
    
    r := strings.NewReplacer("\"", "", "[", "", "]", "", " ", "")
    replaced := r.Replace(str)

    if(strings.Contains(replaced, "," )) {
        return strings.Split(replaced, ",")
    } 
   return []string{replaced}
}