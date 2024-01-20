The analyses-api provides APIs that enable you to perform data analyses on csv data.

Currently there are provisions for:

- mutlivariable linear regression
    

with support for many more analyses operation coming along.

## **Getting started guide**
Read the [postman documentation](https://documenter.getpostman.com/view/17081738/2s9YsT6oZJ)

To start using analyses APIs, you need to -

- You must create an account, and then sign in to the account to obtain an authentication token.
- You must upload data in csv format.
- You must use a valid API Key to send requests to the API analyses endpoints. You can get your API key
- The API only responds to HTTPS-secured communications. Any requests sent via HTTP return an HTTP 301 redirect to the corresponding HTTPS resources.
- The API returns request responses in JSON format. When an API request returns an error, it is sent in the JSON response as an error key.
    

## Authentication

The analyses-api API uses `bearer token` for authentication. You can generate an authentication token by sign in to a created account.

> You must include an API key in each request to the files and analyses endpoints with the `Authorization` request header. 
  

### Authentication error response

If authentication token is missing, malformed, or invalid, you will receive an HTTP 401 Unauthorized response code.
