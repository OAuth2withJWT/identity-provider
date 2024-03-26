# Identity Provider

## Project Setup

### Configure Database Connection:

1. **Create a PostgreSQL Database:** Set up a PostgreSQL database either locally or on a remote server.

2. **Copy the Environment File:** Duplicate the `.example.env` file provided in the root directory of the project and rename it to `.env`.

3. **Configure Environment Variables:** Open the `.env` file in a text editor and fill in the required database connection details, including `DATABASE_NAME`, `DATABASE_USER`, and `DATABASE_PASSWORD`.

### Setup air live-reload

For detailed instructions on setting up Air for live-reload, please follow the [official documentation](https://github.com/cosmtrek/air).

Before starting the Air server, ensure the Go binary directory is added to your system's PATH environment variable. To achieve this, execute:  
```export PATH="$PATH:/home/your_username/go/bin"```
