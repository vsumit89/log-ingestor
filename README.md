## Logswift

A logquery interface which can be used to ingest logs at `http://localhost:3000/` and query them at `http://localhost:3000/query`.

![example-1](./images/example-2.png)
![example-2](./images/example-1.png)

### Demo Video Link

<a href="https://www.loom.com/share/3e3eae3ac9b34e449218872c57256b17?sid=a39e1008-4901-4ab3-b860-2cf37631e901" target="_blank">Demo Video Link</a>

### Techstack used in this project

- NextJS and TailwindCSS - for the client helps in writing less css and faster development
- Golang - for the server which ingests the logs and from which the logs can be queried. Has good support for concurrency which helps in batching some operations.
- Postgres - for storing the logs mainly used as a persistene layer. If meilisearch loses the data, it can be recovered from the postgres database.
- Meilisearch - for indexing the logs and querying them faster. Gives ability for full text search, filtering and typo tolerance.

### Architecture and thought process for the solution

#### Diagram

![architecture img](./images/architecture-diagram-latest.jpeg)

- Ingestion Phase

  1. When the log request comes the server does 2 actions, publishes a task to the rabbitmq and maintains a buffer of max size 1000 in memory. Once the buffer exceed MAX_BUFFER_SIZE (considered 1000) the logs are flushed to postgres. If MAX_BUFFER_SIZE is not reached the logs are flushed to postgres a certain period (5 seconds) in this case.
  2. The task is picked up by the worker and the logs are ingested in meilisearch. This is done in batches of 1000 logs.
  3. For optimising the ingestion into postgres, implemented application level sharding. I am writing it to n different databases. This n is selected based on db config var in config.yaml. This helps in parallelising the writes to the database and hence improving the performance.

- Query Phase

  1. The logs are queried from meilisearch and the results are returned to the client.
  2. The reason for chosing meilisearch is that it provides full text search, typo tolerance and filtering on the data. Since its a search index its optimised for search.
  3. The tradeoff for better search query is that the indexing in meilisearch takes time hence it takes some to reflect the changes in the search results. Until around 50000-100000 the real time search is achieved.

- On Web
  1. Added a debounce of 500ms on the search input to reduce the number of requests to the server.

### Folder structure

Folder structure is described in the README of the respective folders

- [web](web/README.md)
- [server](server/README.md)

## Getting Started

To setup the project locally follow the below steps:

### Prerequisites

- Docker
- Docker-compose

If you want to make any change to the configuration, you can do so by changing the environment variables in the docker-compose file and for server you can update the config.yaml file in the server

### Installation

1. Clone the repo

   ```sh
   git clone https://github.com/dyte-submissions/november-2023-hiring-vsumit89
   ```

2. Run the docker-compose file
   ```sh
    docker-compose up
   ```
3. In case the client fails to start (I was facing issues while testing with docker-compose), you can run the client manually.

   ```sh
     cd web
     npm install
     npm run dev -- -p 3001
   ```

4. The client will be running at `http://localhost:3001` and the server will be running at `http://localhost:3000`
5. If the `logswift-api` and `logswift-consumer` services show `connection refused with rabbitmq`. Restart them using
   `docker-compose restart {service-name}`

## Roadmap
- Ability to add custom MAX_SIZE_BUFFER and FLUSH_INTERVAL for postgres and meilisearch
- Using Clickhouse instead of postgres for storing the logs. Clickhouse has better performance for batch insertion and better data compression.
- Adding unit, integration tests for the server
- Responsiveness of UI
- Adding more worker instances for parallelising the ingestion into meilisearch
- Using rabbitmq for flushing logs into postgres
  with the help of some worker instance
- Add Natural language search too
