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

