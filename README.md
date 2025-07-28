# Software Engineering Task: Ad Bidding Service

## Overview

Ad Bidding service Coupled with Kafka,ClickHouse,prometheus and Graphana

## Task Requirements

Your challenge is to:

1. **Implement ad selection logic** to find winning ads based on various criteria
2. **Implement a tracking endpoint** to record user interactions with ads
3. **Add relevancy scoring** to improve ad matching quality
4. **Implement appropriate validation** for all endpoints
5. **Document your approach** and any assumptions made

## Prerequisites

* [Go](https://golang.org/doc/install)
* [Docker](https://docs.docker.com/engine/install/)
* [Compose](https://docs.docker.com/compose/install/)
* [Kafka](https://kafka.apache.org/)
* [Clickhouse](https://clickhouse.com/)
* [Prometheus](https://prometheus.io/)
* [Graphana](https://grafana.com/)

## Setup & Environment

This repository provides a basic service structure to get started:

```bash
# Build and start the service

docker compose down -v && docker compose up --build -d
# Check service status
curl http://localhost:8080/health
```

## Systems Check
Check if metrics and up and running
```bash
curl -X GET "http://localhost:9090/targets?search="
```
![Prometheus target](https://i.postimg.cc/X79V7Xjv/Screenshot-2025-07-27-at-15-52-10.png)

## Configuration

The service uses environment variables for configuration, using [Kelsey Hightower's envconfig](https://github.com/kelseyhightower/envconfig) library.

Available environment variables:

| Variable        | Description                          | Default |
|-----------------|--------------------------------------|---------|
| APP_NAME        | Application name                     | "Ad Bidding Service" |
| APP_ENVIRONMENT | Running environment                  | "development" |
| APP_LOG_LEVEL   | Log level (debug, info, warn, error) | "info" |
| APP_VERSION     | Application version                  | "1.0.0" |
| SERVER_PORT     | HTTP server port                     | 8080 |
| SERVER_TIMEOUT  | Server timeout for requests          | "30s" |
| BROKER          | Kafka Broker                         | "kafka:9092" |


## Test Setup
Separate cli tool is available to test this
https://github.com/aniruddha-chakraborty/hiring-software-engineer-task-test
Please follow the documentation 

## ✅ Test Run Results
**Ad Selection Test**:
```bash
achakraborty@achakraborty-MacBook-Pro hiring-software-engineer-task-test % go run main.go ad-test                                            
--- Running Simplified & Targeted Ad Logic Tests (First 5 Items) ---

--- [Test 1/5] Testing for: 'Summer Sale Banner' ---
  - Using its own targeting: placement='homepage_top', category='electronics', keyword='summer'
  - Prediction: The API should return 'Summer Sale Banner' in the list of ads.
  - CURL Command: curl -X GET 'http://localhost:8080/api/v1/ads?category=electronics&keyword=summer&limit=4&placement=homepage_top'
  - ACTION: Calling the real API...
  [RESULT] ✅ PASS: Predicted ad 'Summer Sale Banner' was found in the returned list of 4 ads.

--- [Test 2/5] Testing for: 'Winter Clearance Promo' ---
  - Using its own targeting: placement='video_preroll', category='fashion', keyword='clearance'
  - Prediction: The API should return 'Winter Clearance Promo' in the list of ads.
  - CURL Command: curl -X GET 'http://localhost:8080/api/v1/ads?category=fashion&keyword=clearance&limit=4&placement=video_preroll'
  - ACTION: Calling the real API...
  [RESULT] ✅ PASS: Predicted ad 'Winter Clearance Promo' was found in the returned list of 4 ads.

--- [Test 3/5] Testing for: 'Travel Deals Campaign' ---
  - Using its own targeting: placement='article_inline_1', category='travel', keyword='exclusive'
  - Prediction: The API should return 'Travel Deals Campaign' in the list of ads.
  - CURL Command: curl -X GET 'http://localhost:8080/api/v1/ads?category=travel&keyword=exclusive&limit=4&placement=article_inline_1'
  - ACTION: Calling the real API...
  [RESULT] ✅ PASS: Predicted ad 'Travel Deals Campaign' was found in the returned list of 4 ads.

--- [Test 4/5] Testing for: 'Gaming Weekend Blast' ---
  - Using its own targeting: placement='homepage_top', category='gaming', keyword='sale'
  - Prediction: The API should return 'Gaming Weekend Blast' in the list of ads.
  - CURL Command: curl -X GET 'http://localhost:8080/api/v1/ads?category=gaming&keyword=sale&limit=4&placement=homepage_top'
  - ACTION: Calling the real API...
  [RESULT] ✅ PASS: Predicted ad 'Gaming Weekend Blast' was found in the returned list of 4 ads.

--- [Test 5/5] Testing for: 'Home Essentials Discount' ---
  - Using its own targeting: placement='video_preroll', category='home', keyword='deal'
  - Prediction: The API should return 'Home Essentials Discount' in the list of ads.
  - CURL Command: curl -X GET 'http://localhost:8080/api/v1/ads?category=home&keyword=deal&limit=4&placement=video_preroll'
  - ACTION: Calling the real API...
  [RESULT] ✅ PASS: Predicted ad 'Home Essentials Discount' was found in the returned list of 4 ads.

--- Test Summary: 5/5 tests passed. ---
```

**E2E Tracking Test**:

```bash
achakraborty@achakraborty-MacBook-Pro hiring-software-engineer-task-test % go run main.go e2e-tracking-test                                  
--- Running End-to-End Tracking Pipeline Test ---

[PHASE 1] Getting initial row count from ClickHouse...
  - Initial row count in ads_final is: 11394748

[PHASE 2] Sending tracking events to the API...

--- Curl Request #1 ---
curl -X POST 'http://localhost:8080/api/v1/tracking' -H 'Content-Type: application/json' -d '{"event_type":"impression","line_item_id":"li_5c0ebd78-0ca4-4abe-b856-9ab57209ca03","placement":"article_inline_1","user_id":"e2e-user-1","metadata":{"browser":"safari","device":"tablet"}}'

--- Curl Request #2 ---
curl -X POST 'http://localhost:8080/api/v1/tracking' -H 'Content-Type: application/json' -d '{"event_type":"impression","line_item_id":"li_5c0ebd78-0ca4-4abe-b856-9ab57209ca03","placement":"article_inline_1","user_id":"e2e-user-2","metadata":{"browser":"safari","device":"tablet"}}'

... (Requests #3 to #14 omitted for brevity)

--- Curl Request #15 ---
curl -X POST 'http://localhost:8080/api/v1/tracking' -H 'Content-Type: application/json' -d '{"event_type":"impression","line_item_id":"li_5c0ebd78-0ca4-4abe-b856-9ab57209ca03","placement":"article_inline_1","user_id":"e2e-user-15","metadata":{"browser":"safari","device":"tablet"}}'

  - Successfully sent 15 tracking events.

[PHASE 3] Waiting 15 seconds for Kafka and ClickHouse to ingest the data...

[PHASE 4] Querying ClickHouse for final row count...
  - Final row count in ads_final is: 11394763

[RESULT]
  ✅ PASS: Sent 15 events. Row count correctly increased from 11394748 to 11394763.

--- Test Complete ---
```

**Validation tests**:

```bash
achakraborty@achakraborty-MacBook-Pro hiring-software-engineer-task-test % go run main.go validation-test
--- Running Line Item Creation Validation Tests ---

[TEST] Creating a completely valid line item
  [RESULT] ✅ PASS: Received expected status 201

[TEST] Name longer than 100 characters
  [RESULT] ✅ PASS: Received expected status 400

[TEST] Bid less than 0.1
  [RESULT] ✅ PASS: Received expected status 400

[TEST] Budget greater than 10000
  [RESULT] ✅ PASS: Received expected status 400

[TEST] Placement not in the allowed list
  [RESULT] ✅ PASS: Received expected status 400

[TEST] Missing required 'name' field
  [RESULT] ✅ PASS: Received expected status 400

--- Validation Test Summary: 6/6 tests passed. ---
```

## API Structure

The service exposes the following endpoints:

- **POST /api/v1/lineitems**: Create new ad line items with bidding parameters
- **GET /api/v1/ads**: Get winning ads for a specific placement with optional filters (you'll need to implement this)
- **POST /api/v1/tracking**: Record ad interactions (you'll need to implement this)

The complete API specification is available in the OpenAPI document at `api/openapi.yaml`.

## Data Model

The core data model includes:

- **LineItem**: An advertisement with associated bid information
  - `id`: Unique identifier
  - `name`: Display name of the line item
  - `advertiser_id`: ID of the advertiser
  - `bid`: Maximum bid amount (CPM)
  - `budget`: Daily budget for the line item
  - `placement`: Target placement identifier
  - `categories`: List of associated categories
  - `keywords`: List of associated keywords

## Deliverables

Please provide the following:

1. **Ad Selection Logic**: Implement the logic to select winning ads based on placement, categories, and keywords
2. **Tracking Endpoint**: Implement an endpoint to record impressions, clicks, and conversions
3. **Relevancy System**: Develop a scoring mechanism to determine ad relevance
4. **Input Validation**: Add appropriate validation for all API endpoints
5. **Documentation**: Update the README and API docs with your changes

## Evaluation Criteria

Your solution will be evaluated based on:

- **Code quality**: Clean, well-structured, and maintainable code
- **API design**: RESTful design, appropriate error handling, and documentation
- **Implementation quality**: Performance, reliability, and adherence to Go best practices
- **Documentation**: Clear explanation of your approach, design decisions, and trade-offs
- **Testing**: Comprehensive test coverage and consideration of edge cases
- **Innovation**: Creative solutions to the technical challenges presented

## Technical Requirements

- Your solution should be containerized and runnable with docker-compose
- All code should follow Go best practices and conventions
- The API should handle appropriate error cases with meaningful status codes and messages
- Your implementation should consider performance and scaling aspects
- Update the OpenAPI specification to match your implementation

## Storage Solutions

The current implementation uses in-memory storage for simplicity, but this is not suitable for production. You are free to use any storage solution you prefer.
Choose solutions that best fit the requirements and consider factors like scalability, reliability, and performance.

→ Relational database(Postgres/mysql) - for Stroning line items  
→ Column database (druid/clickhouse/bigquery) - for analytics  
→ key value database (keydb/aerospike) - for feature implementation like frequency capping need database like these

## Storage cost savings

```
Using data compression as much as possible, using numerical data is much easier to compress
For example i would not use LineItemID as string, i would use as an integer because when dumping on clickhouse
That single column will weight a hell lot! its not a problem if data is small but we are taking about
billions, if not trillions data points if possible. And we want to show the users as far reports as possible
Even with these measurement we still have detach 1 year old or 6 month old partition and store it somewhere.
And data have to partition based on time.

On the other hand, LineItemID could be a very very useful tool, if we use Lowcardnility() like feature, which
means number of total campaign ID will be storage limit though small sacrifice, it's possible to do data sharding
depends on that LineItemID, we can write a service that can look at the LineItemID string, decrypt it, and figure
out where this data will go. In this way. it's much easier to maintain seperate client with all their rules with
data protection. No one, have access each other data!
```

## Scaling Considerations

As part of your solution, please include a section in your documentation addressing the following questions:

1. How would you scale this service to handle millions of ad requests per minute?
```
→ Stateless microservice architecture
→ Multiple type of database which is easy to scale and can be switch to
  cloud based infra, because when team size is small, it is s better
  to have cloud backup even though it seems like cost is high
→  Async processing of tracking data with kafka and clickhouse
→ All microservice must be able to scale horizontally 
→ When ever writing code, we should think like from time complexity perspective or ops/sec perspective

```
2. What bottlenecks do you anticipate and how would you address them?
```
→ Most complex problem i can imagine is, distribution how much a particulaer bidding server ( there will 100s of bidding server running )
 can spend on campaign. Because we cant forward all the load to one server 
 I would implement some kind algorithm so that a particular bidding server knows how much to spend to prevent over spending
```
3. How would you design the system to ensure high availability and fault tolerance?
```
→ Replication - multiple Kafka brokers, multi-AZ ClickHouse

→ Load balancing - distribute traffic

→ Failover - kubernetes deployment is a must
```
4. What data storage and access patterns would you recommend for different components (line items, tracking events, etc.)?
```
→  For Storing lineitems - Relational database (postgres,mysql/mariadb)
→  For For tracking - Columner database (clickhouse,druid)
```
5. How would you implement caching to improve performance?
```
  I already implemented the type of caching system in this project
```

## Project Structure

```
.
├── api/                    # API documentation and OpenAPI spec
├── cmd/                    # Application entrypoints
│   └── server/             # Main server application
├── internal/               # Private application code
│   ├── config/             # Configuration handling
│   ├── handler/            # HTTP handlers
│   ├── model/              # Data models
│   └── service/            # Business logic
├── docker-compose.yml      # Docker Compose configuration
├── Dockerfile              # Docker build configuration
├── go.mod                  # Go module definition
├── go.sum                  # Go module checksums
└── README.md               # Project documentation
└── .gitignore              # Ignore unimportant files
└── jmx-exporter-config.yml # Kafka metrics exporter config
└── metrics.xml             # Clickhouse config for exporting metrics
└── prometheus.yml          # Prometheus scraping configs are here
└── schema.sql              # Clickhouse database schema
```

## Overall Performance test

```
First: I'm really sorry i didn't use FindMatchingLineItems() function in my getAd()
Because it would drag the performance down i would not be able to practically test
and see what my algorithm did on a performance basis so i wrote that part on the
Add selection algorithm You can take a look there.
```

### Ads selection load test
.lua files are in the https://github.com/aniruddha-chakraborty/hiring-software-engineer-task-test project
```bash
wrk -t10 -c400 -d600s -s ./load_test.lua http://localhost:8080
```

![Cli](https://i.postimg.cc/wM1GF9cP/getad-loadtest-cli.png)
![Graphana](https://i.postimg.cc/MTNhCcTx/getad-loadtest-graphana.png)

### Tracking feature load test
```bash
wrk -t10 -c400 -d300s -s ./tracking_load.lua http://localhost:8080
```
![Cli](https://i.postimg.cc/rwj78Zmf/tracking-load-test-cli.png)
![Graphana](https://i.postimg.cc/dVwzccq7/tracking-load-graphana.png)

