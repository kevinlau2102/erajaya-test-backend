# erajaya-test-backend
Erajaya Technical Test - Backend

## How To Use
1. Clone the repository
  ```bash
  git clone https://github.com/kevinlau2102/erajaya-test-backend.git
  ```
2. Navigate to the project directory:
  ```bash
  cd erajaya-test-backend
  ```
3. Copy the example environment file and configure it:
  ```bash 
  cp .env.example .env
  ```
There are 2 ways to do running
### With Docker
1. Build Docker
  ```bash
  make up
  ```
2. Run Initial UUID V4 for Auto Generate UUID
  ```bash
  make init-uuid
  ```
3. Run Migration and Seeder
  ```bash
  make migrate-seed
  ```

## List API
### GET /products

Get all products with pagination, sorting, and search.

**Query Parameters:**
- `page` (int, optional) - default: 1
- `limit` (int, optional) - default: 10
- `search` (string, optional)
- `sortBy` (string: `name`, `price`, `created_at`)
- `order` (string: `asc`, `desc`)

**Response:**
```json
{
  "products": [...],
  "pagination": {
    "page": 1,
    "limit": 10,
    "count": 100,
    "maxPage": 10
  }
}

You can try the API using this Postman collection: https://api.postman.com/collections/25822863-2feafd6d-19ff-40ce-a9d2-85d9024a6015?access_key=PMAT-01K0DD79WE6XQNMCJ1V119SSSP
