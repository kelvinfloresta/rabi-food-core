# Rabi Food Core

Rabi Food Core is the backend service for the Rabi Food application, responsible for managing restaurants, menus, orders, and user operations within the food delivery platform.

---

## 🚀 Getting Started

### 🧩 Prerequisites

Ensure you have the following installed and configured:

* [Go](https://go.dev/dl/) **v1.25.3**
* [Docker](https://docs.docker.com/get-docker/)
* [Docker Compose](https://docs.docker.com/compose/install/)
* [Task](https://taskfile.dev/installation/)

Install Task with:

```bash
go install github.com/go-task/task/v3/cmd/task@latest
````

---

### ⚙️ Running the Project

Start the development environment with:

```bash
task dev
```

This command will:

* Load environment variables from `.env.test`
* Start the **Go application** (`go run main.go`)
* Launch the **test infrastructure**, including **Postgres** and **pgAdmin**

Once running, you can access:

* **pgAdmin:** [http://localhost:5050](http://localhost:5050) — to browse and verify database data

> The `.env.test` file is already included in the project and contains only non-sensitive values.
> It’s automatically loaded by Task and helps debug environment or connection issues if something fails during startup.

---

### 🧪 Running Tests

To execute the test suite in a clean environment:

```bash
task test
```

This command:

* Starts the **test infrastructure** (Postgres and pgAdmin)
* Clears the Go test cache
* Runs all tests using `gotestsum`

Once running, you can check:

* **pgAdmin:** [http://localhost:5050](http://localhost:5050) — to inspect database results generated during tests

---

### 📊 Observability Mode (Optional)

If you want to test the **observability and logging stack**, use:

```bash
task test-with-logs
```

This command:

* Starts both the **application infrastructure** (Grafana, Prometheus, Loki, Alloy)
  and the **test infrastructure** (Postgres, pgAdmin)
* Runs the tests inside a Dockerized environment
* Streams logs to the observability stack

Once running, you can access:

* **Grafana:** [http://localhost:3100](http://localhost:3100) — to view application and test logs
* **pgAdmin:** [http://localhost:5050](http://localhost:5050) — to verify database state

> This mode is intended for validating **log ingestion and observability configuration**, not for everyday development.

---

### 🧰 Additional Utilities

* **Generate mocks using Mockery**

  ```bash
  task mockgen
  ```

* **Validate Alloy configuration**

  ```bash
  task validate_alloy
  ```

* **Clean up Docker resources**

  ```bash
  task clean_docker
  ```

---

## 🧭 `app/` Directory Structure

The `app/` folder contains the core logic of the application.
It is organized to keep the codebase **modular**, **readable**, and **easy to maintain**, following principles inspired by **Clean Architecture** and **Domain-Driven Design (DDD)**.
Each subdirectory has a specific responsibility and plays a distinct role in the system design.

---

## 📂 Directory Overview

### **`domain/`**

Defines the **business entities** and **core rules** of the system.
Represents the application's domain model, independent from frameworks or infrastructure.

---

### **`usecases/`**

Contains the **application logic** that coordinates business operations.
Each folder (e.g., `user_case`, `product_case`, `tenant_case`) represents a **business context** or **functional area**.

---

### **`libs/`**

Holds **infrastructure-related libraries and adapters** that connect the application to external systems.

Includes:

* `database/` – gateways, adapters, and ORM integrations
* `http/` – controllers, routes, and web adapters
* `logger/` – centralized logging utilities
* `validator/` – input validation helpers
* `di/` – dependency injection setup

This directory may include both **interfaces** and **implementations** used by the use cases.

---

### **`config/`**

Manages environment configuration for development, testing, and production.

---

### **`app_context/`**

Defines shared runtime context — such as session handling or dependency resolution.

---

### **`fixtures/`**

Provides reusable components for tests, including data builders, mocks, and helpers.

---

## 🧠 Summary

```
app/
├── app_context/     # Shared runtime context
├── config/          # Environment configuration
├── domain/          # Core business entities and rules
├── fixtures/        # Test data and mocks
├── libs/            # Infrastructure and external integrations
│   ├── database/    # Gateways and adapters
│   ├── http/        # Controllers and routes
│   ├── logger/      # Logging utilities
│   ├── validator/   # Input validation
│   └── di/          # Dependency injection setup
└── usecases/        # Business workflows and application logic
```

---

## 💡 Design Rationale

This organization helps to:

* Keep **business logic isolated** from infrastructure concerns
* Allow **easy replacement or extension** of technical components (database, web framework, etc.)
* Improve **readability**, **testability**, and **consistency** across the codebase
* Maintain a clear structure as the project evolves
