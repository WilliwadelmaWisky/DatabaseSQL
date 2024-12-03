# SQL Database
An extremely simple local sql database that supports only a few of the sql keywords.
The data can be requested/updated from the server by http post request at root.
The data is send in json format as http response if the request is in correct format.

## Getting Started
**Go must be installed!**
- Download the application from releases
- Run `go run .`, the server starts at port 9000

## Supported SQL
**CREATE**
```sql
CREATE TABLE <table> (
    <name> INT
    <name> VARCHAR(<len>)
)
```

**SELECT**
```sql
SELECT <column> FROM <table>
WHERE <condition>
ORDER BY <column> <ASC|DESC>
```

**UPDATE**
```sql
UPDATE <table>
SET (
    <name> = 10,
    <name> = 'Text'
)
WHERE <condition>
```

**DELETE**
```sql
DELETE <table>
WHERE <condition>
```
