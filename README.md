# go-api
API for customer orders

DATABASE:
- Run Postgres in local 
- Run the scripts `db_scripts.sql` and `insert_scripts.sql` for data setup

Run in local:
- Clone the repo in the local.
- add Postgres DB configurations in initDB func, inside main.go file. 
- Then run the command `go run main.go'.

API:
Run the api with the header token `Token` and value `hunter2` to authorize the api.
- GET    /health                   : To check the api status.
- GET    /api/v1/orders            : List all the orders.
- GET    /api/v1/orders/:dateRange : Get orders between order date range. 

- POST /api/v1/data/:fileName      : load data into Postgres DB respective table.  
