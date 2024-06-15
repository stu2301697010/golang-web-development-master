# Book Management Web API

## Description

This project is a web API for managing a collection of books, built with GoLang and the Gin web framework. The API allows users to create, read, update, and delete book records in a PostgreSQL database. The project is containerized using Docker and Docker Compose, making it easy to deploy.

Each book in the collection has the following attributes:
- Title (string): The title of the book.
- ISBN (string): The International Standard Book Number, a unique identifier for books.
- Author (string): The author of the book.
- Year of Publication (integer): The year the book was published.

## Endpoints

- `GET /books`: Retrieve a collection of all books with pagination.
- `GET /books/:id`: Retrieve details of a single book by its identifier.
- `POST /books`: Create a new book record.
- `PUT /books/:id`: Update the information of an existing book.
- `DELETE /books/:id`: Delete a book specified by its identifier.

## Getting Started

### Prerequisites

- Docker
- Docker Compose

### Installation

1. Clone the repository:

```bash
git clone https://github.com/stu2301697010/golang-web-development-master.git
cd golang-web-development-master/
```
Build and run the application using Docker Compose:

```bash
docker-compose up --build
```
This will start the application on http://localhost:8080.

Testing with cURL
You can test the API endpoints using cURL commands as shown below.

1. Retrieve all books with pagination
```bash
curl -X GET "http://localhost:8080/books?page=1&limit=10" -H "accept: application/json"
```

2. Retrieve details of a single book by its ID
```bash
curl -X GET "http://localhost:8080/books/1" -H "accept: application/json"
```

3. Create a new book record
```bash
curl -X POST "http://localhost:8080/books" -H "accept: application/json" -H "Content-Type: application/json" -d '{
  "title": "New Book Title",
  "isbn": "1122334455",
  "author": "New Author",
  "year": 2022
}'
```

4. Update the information of an existing book
```bash
curl -X PUT "http://localhost:8080/books/1" -H "accept: application/json" -H "Content-Type: application/json" -d '{
  "title": "Updated Book Title",
  "isbn": "1122334455",
  "author": "Updated Author",
  "year": 2023
}'
```

5. Delete a book by its ID
```bash
curl -X DELETE "http://localhost:8080/books/1" -H "accept: application/json"
```
**Environment Variables**
The following environment variables are used to configure the PostgreSQL database:

- `DB_HOST`: The hostname of the PostgreSQL database (default: db).
- `DB_USER`: The PostgreSQL user (default: youruser).
- `DB_PASSWORD`: The PostgreSQL user's password (default: yourpassword).
- `DB_NAME`: The name of the PostgreSQL database (default: yourdb).
- `DB_PORT`: The port on which PostgreSQL is running (default: 5432).

These variables are defined in the docker-compose.yml file.

## Load Sample Data
The application will load sample data from a CSV file (books.csv) into the PostgreSQL database when it starts. The sample data consists of the first 1000 records from the Books Dataset on Kaggle. Ensure the CSV file is placed in the project directory.

## Docker Compose Configuration
The `docker-compose.yml` file defines two services:

- `db`: The PostgreSQL database service.
- `web`: The GoLang web API service.

Volumes are used to persist the PostgreSQL data.