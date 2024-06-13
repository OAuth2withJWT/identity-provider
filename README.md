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

### Database migrations

We use ```golang-migrate/migrate``` for managing database migrations. To create and apply migrations, follow these steps:

1. **Install migrate CLI:** Install the `migrate` command-line tool by following the instructions provided in the [official repository](https://github.com/golang-migrate/migrate/tree/master).

2.  **Create a Migration:** Run the following command to create a new migration file:
    
    ```make create_migration NAME=<migration_name>```

3.  **Apply Migrations:** To apply pending migrations and update the database schema, run:
    
    ```make migration_up [N=<number_of_migrations_to_apply>]```

    
4.  **Rollback Migrations:** To rollback the applied migrations, run:
    
    ```make migration_down [N=<number_of_migrations_to_rollback>]```
    

### Using Redis with go-redis/v9

We use ```go-redis/v9``` as our in-memory database solution. To integrate it into your project, follow these steps:

1. **Install Redis**

    ```go get github.com/redis/go-redis/v9```

2. **Start Redis server**: To start the Redis server, simply run the following command in your terminal:

    ```redis-server```