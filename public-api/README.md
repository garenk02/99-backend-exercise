# Public API

These are the public facing APIs that can be called by external clients such as mobile applications or the user facing website.

### Setup
To get the `public api` up and running, you will need **Golang 1.23+** version.

### Install Golang
Head over to https://go.dev/dl then choose based on your operating system.

Once you have golang set up, then we can proceed to run our web applications:

```bash
# Copy .env.example to .env
cp .env.example .env

# Install the dependencies
go mod tidy

# Run the user service
go run *.go
```

You can adjust the value on the `.env` file depend on your needs.

#### Get listings
Get all the listings available in the system (sorted in descending order of creation date). Callers can use `page_num` and `page_size` to paginate through all the listings available. Optionally, you can specify a `user_id` to only retrieve listings created by that user.

```
URL: GET /public-api/listings

Parameters:
page_num = int # Default = 1
page_size = int # Default = 10
user_id = str # Optional
```
```json
{
    "result": true,
    "listings": [
        {
            "id": 1,
            "listing_type": "rent",
            "price": 6000,
            "created_at": 1475820997000000,
            "updated_at": 1475820997000000,
            "user": {
                "id": 1,
                "name": "Suresh Subramaniam",
                "created_at": 1475820997000000,
                "updated_at": 1475820997000000,
            },
        }
    ]
}

```

#### Create user
```
URL: POST /public-api/users
Content-Type: application/json
```
```json
Request body: (JSON body)
{
    "name": "Lorel Ipsum"
}
```
```json
Response:
{
    "user": {
        "id": 1,
        "name": "Lorel Ipsum",
        "created_at": 1475820997000000,
        "updated_at": 1475820997000000,
    }
}
```

#### Create listing
```
URL: POST /public-api/listings
Content-Type: application/json
```
```json
Request body: (JSON body)
{
    "user_id": 1,
    "listing_type": "rent",
    "price": 6000
}
```
```json
Response:
{
    "listing": {
        "id": 143,
        "user_id": 1,
        "listing_type": "rent",
        "price": 6000,
        "created_at": 1475820997000000,
        "updated_at": 1475820997000000,
    }
}
```
