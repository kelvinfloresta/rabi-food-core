# Rabi Food Core
Rabi Food Core is the backend service for the Rabi Food application, which manages food delivery operations including restaurants, menus, orders, and user management.


## 🚀 Getting Started

### 🧩 Prerequisites

Make sure you have the following installed and configured on your machine:

* [Go](https://go.dev/dl/) **v1.25.3**
* [Docker](https://docs.docker.com/get-docker/)
* [Docker Compose](https://docs.docker.com/compose/install/)
* [Taskfile](https://taskfile.dev/installation/)

You can install Taskfile with:

```bash
go install github.com/go-task/task/v3/cmd/task@latest
```

---

### ⚙️ Running the Project

To start the development environment, simply run:

```bash
task dev
```

This command will:

* Load environment variables from `.env.test`
* Start the Go application (`go run main.go`)
* Spin up the **test infrastructure**, including **Postgres** and **pgAdmin**

Once running, you can access:

* **pgAdmin:** [http://localhost:5050](http://localhost:5050) — to browse and verify database data
> The `.env.test` file is already included in the project and contains only non-sensitive values.
> It’s used automatically when running the tasks, which helps when debugging environment or connection issues.


---

### 🧪 Running Tests

To run the full test suite (with Postgres and pgAdmin containers automatically started):

```bash
task test
```

This command:
* Starts **test infrastructure** (Postgres, pgAdmin)
* Runs the tests inside a Dockerized environment
* Streams logs to the observability stack

Once running, you can access:

* **pgAdmin:** [http://localhost:5050](http://localhost:5050) — to check 

If you need to **test the observability stack and verify log ingestion** into Grafana or Loki, use:

```bash
task test-with-logs
```

This command:

* Starts both the **application infrastructure** (Grafana, Prometheus, Loki, Alloy)
  and the **test infrastructure** (Postgres, pgAdmin)

Once running, you can access:
* **Grafana:** [http://localhost:3100](http://localhost:3100) — to inspect application and test logs

This mode is mainly used to **validate the log pipeline and observability configuration**, not for regular development.

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
