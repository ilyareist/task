# API Documentation

## Table of Contents

<!-- TOC depthFrom:2 depthTo:6 updateOnSave:true withLinks:true -->

- [Table of Contents](#table-of-contents)
- [Accounts Collection `/api/accounts/v1/accounts`](#accounts-collection-apiaccountsv1accounts)
    - [List All Accounts](#list-all-accounts)
        - [Request](#request)
        - [Responses](#responses)
            - [Success response](#success-response)
    - [Create a New Account](#create-a-new-account)
        - [Request](#request-1)
        - [Responses](#responses-1)
            - [Success response](#success-response-1)
            - [Error responses](#error-responses)
                - [406 Not Acceptable](#406-not-acceptable)
                - [500 Internal Server Error](#500-internal-server-error)
- [Account `/api/accounts/v1/accounts/{accountid}`](#account-apiaccountsv1accountsaccountid)
    - [Get account by ID](#get-account-by-id)
        - [Request](#request-2)
        - [Responses](#responses-2)
            - [Success response](#success-response-2)
            - [Error responses](#error-responses-1)
                - [404 Not Found](#404-not-found)
    - [Delete account by ID](#delete-account-by-id)
        - [Request](#request-3)
        - [Responses](#responses-3)
            - [Success response](#success-response-3)
            - [Error responses](#error-responses-2)
                - [404 Not Found](#404-not-found-1)
- [Payments Collection `/api/payments/v1/payments`](#payments-collection-apipaymentsv1payments)
    - [List All Payments](#list-all-payments)
        - [Request](#request-4)
        - [Responses](#responses-4)
            - [Success response](#success-response-4)
    - [Create a New Payment](#create-a-new-payment)
        - [Request](#request-5)
        - [Responses](#responses-5)
            - [Success response](#success-response-5)
            - [Error responses](#error-responses-3)
                - [400 Bad Request](#400-bad-request)
                - [404 Not Found](#404-not-found-2)
                - [406 Not Acceptable](#406-not-acceptable-1)
                - [500 Internal Server Error](#500-internal-server-error-1)
- [Payments by Account `/api/payments/v1/payments/{accountid}`](#payments-by-account-apipaymentsv1paymentsaccountid)
    - [Get Payments for Account](#get-payments-for-account)
        - [Request](#request-6)
        - [Responses](#responses-6)
            - [Success response](#success-response-6)
            - [Error responses](#error-responses-4)
                - [500 Internal Server Error](#500-internal-server-error-2)

<!-- /TOC -->

## Accounts Collection `/api/accounts/v1/accounts`

### List All Accounts

Returns all accounts registered in the system.

#### Request

**URL**: `/api/accounts/v1/accounts`  
**Method**: `GET`

```bash
curl --include \
'http://0.0.0.0:8099/api/accounts/v1/accounts'
```

#### Responses

##### Success response

**HTTP Status**: `200 OK`

```json
[
  {
    "id": "alice456",
    "balance": 999.99,
    "currency": "USD"
  },
  {
    "id": "bob123",
    "balance": 87.78,
    "currency": "USD"
  }
]
```

### Create a New Account

You may create new account using this action. It takes a JSON object containing an id, initial balance and currency.

#### Request

**URL**: `/api/accounts/v1/accounts`  
**Method**: `POST`

```bash
curl --include \
     --request POST \
     --header "Content-Type: application/json" \
     --data-binary "{
    \"id\": \"john789\",
    \"balance\": 55.00,
    \"currency\": \"USD\"
}" \
'http://0.0.0.0:8099/api/accounts/v1/accounts'
```

#### Responses

##### Success response

**HTTP Status**: `200 OK`

```json
{}
```

##### Error responses

###### 406 Not Acceptable

**Condition**: If validation of incoming payload not passed.  
**HTTP Status**: `406 Not Acceptable`

```json
{
  "error": "validation error message"
}
```

###### 500 Internal Server Error 

**HTTP Status**: `500 Internal Server Error`

```json
{
  "error": "error message"
}
```

## Account `/api/accounts/v1/accounts/{account_id}`

### Get account by ID

Returns a read model of an account.

#### Request

**URL**: `/api/accounts/v1/accounts/{account_id}`  
**Method**: `GET`  
**Parameters**:
  - `account_id` - _string_ -- ID of the Account in the form of an alphanumeric string [a-zA-Z0-9].

```bash
curl --include \
'http://0.0.0.0:8099/api/accounts/v1/accounts/bob123'
```

#### Responses

##### Success response

**HTTP Status**: `200 OK`

```json
{
  "account": 
  {
    "id": "bob123",
    "balance": 87.78,
    "currency": "USD"
  }
}
```

##### Error responses

###### 404 Not Found 

**Condition**: If specified account not found.  
**HTTP Status**: `404 Not Found`

```json
{
  "error": "unknown account"
}
```

### Delete account by ID

This method uses to delete account from the system. Actually mark it as deleted.

#### Request

**URL**: `/api/accounts/v1/accounts/{account_id}`  
**Method**: `DELETE`  
**Parameters**:
  - `account_id` - _string_ -- ID of the Account in the form of an alphanumeric string [a-zA-Z0-9].

```bash
curl --include \
     --request DELETE \
'http://0.0.0.0:8099/api/accounts/v1/accounts/test1'
```

#### Responses

##### Success response

**HTTP Status**: `200 OK`

```json
{}
```

##### Error responses

###### 404 Not Found 

**Condition**: If specified account not found.  
**HTTP Status**: `404 Not Found`

```json
{
  "error": "unknown account"
}
```

## Payments Collection `/api/payments/v1/payments`

### List All Payments

Returns all payments, registered in the system.

#### Request

**URL**: `/api/payments/v1/payments`  
**Method**: `GET`

```bash
curl --include \
'http://0.0.0.0:8099/api/payments/v1/payments'
```

#### Responses

##### Success response

**HTTP Status**: `200 OK`

```json
[
  {
    "account": "bob123",
    "amount": 12.34,
    "to_account": "alice456",
    "direction": "outgoing"
  },
  {
    "account": "alice456",
    "amount": 12.34,
    "from_account": "bob123",
    "direction": "incoming"
  }
]
```

### Create a New Payment

You may create new payment using this action. It takes a JSON object containing an from [account ID], amount of transferring money and to [account ID].

#### Request

**URL**: `/api/payments/v1/payments`  
**Method**: `POST`

```bash
curl --include \
     --request POST \
     --header "Content-Type: application/json" \
     --data-binary "{
    \"from\": \"bob123\",
    \"amount\": 12.34,
    \"to\": \"alice456\"
}" \
'http://0.0.0.0:8099/api/payments/v1/payments'
```

#### Responses

##### Success response

**HTTP Status**: `200 OK`

```json
{}
```

##### Error responses

###### 400 Bad Request 

**Condition**: If source account doesn't have enough money for transfer.  
**HTTP Status**: `400 Bad Request`

```json
{
  "error": "insufficient money on source account"
}
```

###### 404 Not Found 

**Condition**: If source or target account not found. Details in error message.  
**HTTP Status**: `404 Not Found`

```json
{
  "error": "unknown account"
}
```

###### 406 Not Acceptable

**Condition**: If validation of incoming payload not passed.  
**HTTP Status**: `406 Not Acceptable`

```json
{
  "error": "validation error message"
}
```

###### 500 Internal Server Error 

**HTTP Status**: `500 Internal Server Error`

```json
{
  "error": "error message"
}
```

## Payments by Account `/api/payments/v1/payments/{account_id}`

### Get Payments for Account

Returns payments list for an account.

#### Request

**URL**: `/api/payments/v1/payments/{account_id}`  
**Method**: `GET`
**Parameters**:
  - `account_id` - _string_ -- ID of the Account in the form of an alphanumeric string [a-zA-Z0-9].

```bash
curl --include \
'http://0.0.0.0:8099/api/payments/v1/payments/bob123
'
```

#### Responses

##### Success response

**HTTP Status**: `200 OK`

```json
[
  {
    "account": "bob123",
    "amount": 12.34,
    "to_account": "alice456",
    "direction": "outgoing"
  }
]
```

##### Error responses

###### 500 Internal Server Error 

**HTTP Status**: `500 Internal Server Error`

```json
{
  "error": "error message"
}
```
