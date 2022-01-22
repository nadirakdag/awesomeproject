
# Awesome Project

This project has two end points in-memory and records. It has api prefix in the url and  handler integration tests. Also supports docker. Also this repo has github action every commit sended to master automatically deploy to heroku.

## Heroku address

[https://awesomeproject-nadir.herokuapp.com](https://awesomeproject-nadir.herokuapp.com)

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

## Run With Docker

Clone the project

```bash
  git clone https://github.com/nadirakdag/awesomeproject
```

Go to the project directory

```bash
  cd awesomeproject
```

Build docker image

```bash
  docker build . -t awesomeproject:latest --build-arg PORT=9090
```

Run docker image 

```bash
 docker run -d -p 9090:9090 awesomeproject:latest
```

  
## Running Tests

To run tests, run the following command

```bash
  go test ./...
```

## ENV Variables
| Name | Default Value     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `PORT` | `8080` | if found use the value of for server port |

  
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

Example curl request

````bash
curl -v -d '{"startDate": "2016-04-01", "endDate":"2016-04-20", "minCount":5000, "maxCount": 6000}' -X POST  https://awesomeproject-nadir.herokuapp.com/api/records
````


#### Get In-Memory

```http
  GET /api/in-memory
```

Returns all of data key values

Example curl request

````bash
curl -v -X GET http://awesomeproject-nadir.herokuapp.com/api/in-memory
````

#### Get In-Memory by Key

```http
  GET /api/in-memory?key={key}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `key`      | `string` | **Required**. key of item to fetch |

Example curl request

````bash
curl -v -X GET http://awesomeproject-nadir.herokuapp.com/api/in-memory?key=test
````


#### Create in-memory key value pair

```http
  POST /api/in-memory
```

| Body | Type     | Description                       | Location |
| :-------- | :------- | :-------------------------------- | :------ |
| `KeyValuePair Object` | `Json` | **Required**. Key Valur information to create | Body |


KeyValuePair Object
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
Example curl request

````bash
curl -v -X POST -d '{"key":"test", "value":"test-value"}' http://awesomeproject-nadir.herokuapp.com/api/in-memory
````