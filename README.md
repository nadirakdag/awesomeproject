
# Awesome Project

## Run Locally

Clone the project

```bash
  git clone https://github.com/nadirakdag/awesomeproject
```

Go to the project directory

```bash
  cd awesomeproject
```

Install dependencies

```bash
  go get
```

Start the server

```bash
  go run main.go
```
  
## Running Tests

To run tests, run the following command

```bash
  go test
```

  
## API Reference

#### Get Records

```http
  POST /api/records
```

| Body | Type     | Description                       | Location |
| :-------- | :------- | :-------------------------------- | :------ |
| `Record Filter Object` | `Json` | **Required**. Records Filter object for filtering records | Body |

##### RecordFilter Object
| Body | Type     | Description                       | Format |
| :-------- | :------- | :-------------------------------- | :------ |
| `startDate` | `string` | **Required**. Start Date for Filtering | YYYY-MM-DD |
| `endDate` | `string` | **Required**. End Date for Filtering | YYYY-MM-DD |
| `minCount` | `int` | **Required**. Minimum Count for Filtering | |
| `maxCount` | `int` | **Required**. Maximum Count for Filtering | |

Example
````Json
{
    "startDate": "2016-01-26",
    "endDate": "2018-02-02",
    "minCount": 2700,
    "maxCount": 3000
}
````

#### Get In-Memory

```http
  GET /api/in-memory
```

Returns all of data key values

#### Get In-Memory by Key

```http
  GET /api/in-memory?key={key}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `key`      | `string` | **Required**. key of item to fetch |


#### Create Bill

```http
  POST /api/in-memory
```

| Body | Type     | Description                       | Location |
| :-------- | :------- | :-------------------------------- | :------ |
| `InMemory Object` | `Json` | **Required**. Key Valur information to create | Body |


##### InMemory Object
| Body | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `key` | `string` | **Required**. Key |
| `value` | `string` | **Required**. Value |

Example
````Json
{
    "key": "info",
    "value": "nadirakdag"
}
````
