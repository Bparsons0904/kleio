# Database Migration for Production

This guide explains how to move your existing Kleio SQLite database from the local development environment to the production NAS storage location.

## Prerequisites

- The Kleio service should be stopped.
- The NAS mount at `/mnt/nas-direct` must be available.
- The production data directory `/mnt/nas-direct/kleio_data` must be created.

## Migration Steps

1.  **Stop the Existing Container**

    If your local Kleio container is running, stop it to ensure the database file is not in use.

    ```bash
    docker stop kleio
    ```

2.  **Copy the Database File**

    The local development setup stores the database in `~/kleio_data`. Copy the `kleio.db` file from this directory to the new NAS location.

    ```bash
    cp ~/kleio_data/kleio.db /mnt/nas-direct/kleio_data/
    ```

3.  **Set Permissions**

    Ensure the user running the Docker container (typically UID 1000) has ownership of the database file on the NAS.

    ```bash
    sudo chown 1000:1000 /mnt/nas-direct/kleio_data/kleio.db
    ```

4.  **Start the Production Service**

    The `docker-compose.prod.yml` file is already configured to use the new database path on the NAS. You can now start the service using the production compose file.

    ```bash
    docker-compose -f docker-compose.prod.yml up -d
    ```

The Kleio service will now run using the migrated database from your NAS.
