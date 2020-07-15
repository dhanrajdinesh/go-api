# go-api

API to get Customer Orders Details.

DATABASE:

- Run Postgres in local 
- Run the scripts `db_scripts.sql` and `insert_scripts.sql` for data setup

Run in local:

- Clone the repository in to local directory.
- Update Postgres DB configurations in "initDB func", inside main.go file. 
- To run the application, runt the command `go run main.go' from local poroject directory.

API:

Run the api with the header token `Token` and value `hunter2` to authorize the api.

- GET    /health                   : To check the api status.
- GET    /api/v1/orders            : List all the orders.
- GET    /api/v1/orders/:dateRange : Get orders between order date range. 

- POST /api/v1/data/:fileName      : It will load data into Postgres DB.  
