# Description
The SmartCar API acts as a forward proxy, standardizing the responses from the GM API. The API documentation is hosted at:
https://app.swaggerhub.com/apis/pink-cupcake/SmartCar/1.01

# To run the app_api service
Requires go 1.13
```bash
cd SmartCar
go run app_api
```
Note: the executible binary is included and can be run directly. If it fails - check if the environment variables were set.

# To generate Swagger documentation
```bash
go generate
```

# To test the API
```bash
go test
```