# GO-GIN-SWAGGER TEST

A REST API microservice using Gin with Zap, Open Telemetry, and Mongo DB?In Memory DB Integrated.

## TODO 

- Decoupling of the logger and Otel API from the App structure
- Common identifier for fetching data on any of the Repository interface implementations
- Make a connection to a Otel provider (AWS, Uptrace, Jaeger, etc..)
- Dockerize the microservice 
- Implement a CLI entrypoint trough scripting (Task, bash or make)
- Use of kubernetes to deploy the  microservcice and its dependencies
- De-hardwire persistence and host information from code
- Find a way to make a better connection for the mongo db client
- Decouple the in memory DB, using a cache tool (e.g redis)
- aggregate live reload
- aggregate session management 
- implement graceful restart and stop as https://gin-gonic.com/docs/examples/graceful-restart-or-stop/
- Readiness and Liveness probes