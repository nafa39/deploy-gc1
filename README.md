[![Review Assignment Due Date](https://classroom.github.com/assets/deadline-readme-button-24ddc0f5d75046c5622901739e7c5dd533143b0c8e959d652212380cedb1ea36.svg)](https://classroom.github.com/a/GMrD03Jz)
# Graded Challenge 1 - P3

Graded Challenge ini dibuat guna mengevaluasi pembelajaran pada Hacktiv8 Program Fulltime Golang khususnya pada pembelajaran mongoDB dan implementasi terhadap aplikasi golang
# **Assignment: Building a Cloud-Echo CRUD Application with Microservices, Docker, GCP, and MongoDB**

## **1. Assignment Objectives**
Graded Challenge 1 ini dibuat guna mengevaluasi pemahaman mongoDB sebagai berikut:

- Mampu memahami Mongo DB
- Mampu Mengimplementasikan konsep Microservices
- Mampu mengimplementasikan Docker 
- Mampu mengimplementasikan Cron Job 

## **2. Assignment Directions**
1. **Environment Setup**:
   - **GCP**: Set up a project in Google Cloud Platform (GCP) for Cloud Run.
   - **Docker**: Install Docker and configure your local development environment.

2. **Microservices Development**:
   - **User Service**: 
     - Endpoints: `POST /users`, `GET /users/:id`, `PUT /users/:id`, `DELETE /users/:id`
   - **Product Service**: 
     - Endpoints: `POST /products`, `GET /products/:id`, `PUT /products/:id`, `DELETE /products/:id`
   - **Order Service**: 
     - Endpoints: `POST /orders`, `GET /orders/:id`, `PUT /orders/:id`, `DELETE /orders/:id`
     - **Cron Job**: Implement a cron job functionality within the Order Service endpoint. This should be triggered via a scheduled HTTP request (using an endpoint like `GET /orders/update-status` that can be scheduled to run daily).

3. **Database Integration**:
   - **MongoDB Setup**:
     - Design MongoDB collections for Users, Products, and Orders.
     - Integrate each microservice with its respective MongoDB collection.

4. **Containerization**:
   - **Dockerize Microservices**: 
     - Write Dockerfiles for each microservice.
     - Build and tag Docker images.

5. **Deployment**:
   - **Cloud Run Deployment**:
     - Upload the Docker images to Google Container Registry (GCR).
     - Deploy each Docker image to Cloud Run.
     - Configure environment variables (e.g., MongoDB connection strings) directly in Cloud Run settings.
   - **Cron Job Scheduling**:
     - Use Google Cloud Scheduler to trigger the cron job endpoint (`GET /orders/update-status`) daily.



## **3. Database Schema**
### **User Service**
| Field        | Type   | Description                 |
|--------------|--------|-----------------------------|
| `id`         | String | Unique identifier for user  |
| `name`       | String | Name of the user            |
| `email`      | String | Email address of the user   |
| `created_at` | Date   | Timestamp when the user was created |

### **Product Service**
| Field        | Type   | Description                 |
|--------------|--------|-----------------------------|
| `id`         | String | Unique identifier for product |
| `name`       | String | Name of the product         |
| `price`      | Float  | Price of the product        |
| `stock`      | Int    | Quantity in stock           |
| `created_at` | Date   | Timestamp when the product was created |

### **Order Service**
| Field        | Type   | Description                 |
|--------------|--------|-----------------------------|
| `id`         | String | Unique identifier for order |
| `user_id`    | String | ID of the user who placed the order |
| `product_id` | String | ID of the product being ordered |
| `quantity`   | Int    | Quantity of the product ordered |
| `total`      | Float  | Total price of the order    |
| `created_at` | Date   | Timestamp when the order was created |

## **4. Expected Results**
- **Microservices**: Each microservice should be independently deployable and capable of handling its own CRUD operations.
- **API Documentation**: Provide API documentation for each microservice, detailing the endpoints and expected input/output.
- **Database Operations**: CRUD operations should be fully functional and correctly interact with MongoDB.
- **Deployment**: The application should be successfully deployed and accessible via GCP.


###  Assignment Submission

Push Assigment yang telah Anda buat ke akun Github Classroom Anda masing-masing.

----------

## Assignment Rubrics

Aspect : 
|Criteria|Meet Expectations|Points|
|---|---|---|
|Problem Solving|5 API Endpoints are implemented and working correctly (@15 each) |75 pts |
|Database Design |MongoDB database meets the required specifications |10 pts|
||Database queries are efficient and appropriately indexed |5 pts|
|Readability|Code is well-documented and easy to read |5 pts|
||Code includes appropriate comments and documentation |5 pts|


Notes:
- Don't rush through the problem or try to solve it all at once.
- Don't copy and paste code from external sources without fully understanding how it works.
- Don't hardcode values or rely on assumptions that may not hold true in all 
cases.
- Don't forget to handle error cases or edge cases, such as invalid input or unexpected behavior.
- Don't hesitate to refactor your code or make improvements based on feedback or new insights.



Total Points : 100

Notes Deadline : W2D1 - 9.00AM

Keterlambatan pengumpulan tugas mengakibatkan skor GC 1 menjadi 0.

