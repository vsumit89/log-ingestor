## Folder structure for the server

- cmd: contains api service and the main.go(entrypoint) file corresponding to it
- docker: contains the dockerfile for the api service
- internal: contains all the file related to this particular api service

  - app: contains files related to app config and app initialization
  - db : contains files related to storage, repository implementation

    - meilisearch: contains meilisearch implementation of search index interface
    - postgres: contains postgres implementation of db interface

  - domain: contains interfaces which define the core business logic
  - dtos: contains the data transfer objects
  - http: contains the http handler and router
  - models: contains the models for the database
  - repository: contains the repository interface which defines the db operations for the models

- pkg : contains the common packages used across the project

- config.yaml: contains the configuration for the server
