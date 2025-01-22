# Order Management Service

## Directory Structure 
```
oms-service/
├── bulkOrderFailures.txt
├── docker-compose.yml
├── go.mod
├── go.sum
├── handler/
│   ├── allHandlers.go
│   ├── bulkOrder.go
│   └── rabbitmq-send-recieve.go
├── main.go
└── models/
    └── order.go
```

## Completion Status
- [x] Single order creation
- [x] Bulk order creation
- [x] Using Message Queue to send Bulk Messages
- [x] Implement Worker Pool for efficient tasks distribution
- [ ] Websockets Implementation

## Project Setup

1. Setup Dependencies
``` golang
  go get "github.com/lib/pq" "github.com/gorilla/mux" "github.com/rabbitmq/amqp091-go"
```
2. Initialize RabbitMQ Image
``` golang
  docker compose up
```
3. Run Main Program
``` golang
  go run main.go
```
4. Make Database: ``` orderDB ``` in PostgreSQL

## API Endpoints

### 1. **Validate and Create Order**
   - Path: `http://localhost:9000/order`
   - Method: `POST`
   - Description: Initiates bulk order processing through RabbitMQ
   - Request Body:
   ```sh
     {
      "product_id": 2,
      "customer_id": 2
     }
   ```
   ![image](https://github.com/user-attachments/assets/1b7b7e9c-b59f-42fb-8918-dfc3a638b1ec)


###  2. **Validate and Create Bulk Orders**
   - Path: `http://localhost:9000/bulkorder`
   - Method: `POST`
   - Description: Initiates bulk order processing through RabbitMQ
   - Request Body:
   ```sh
       {
      	{2, 1},
      	{6, 7},
      	{2, 2},
      	{3, 2},
      	{4, 1},
      	{9, 9},
      }
   ```
  ![image](https://github.com/user-attachments/assets/5a42a64a-647b-49e3-b599-5c1342352fa1)

### Logging Failed Bulk Requests
![image](https://github.com/user-attachments/assets/296c0038-d932-43fa-93ac-377a3952b986)

### Successful Order Creation Email
![image](https://github.com/user-attachments/assets/775658fd-99af-47d3-8928-757ee28f3c2f)


## Ideation
![IMG_20250120_205847](https://github.com/user-attachments/assets/8b3f6229-79bc-4a83-9e97-1ef85e645d39)
