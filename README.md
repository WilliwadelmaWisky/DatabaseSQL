# Database SQL

<p align="justify">
    An extremely simple local sql database that supports only a few of the sql keywords. The data can be requested/updated from the server by http post request at root. The data is send in json format as http response if the request sql is in correct format. The application supports only a single sql query send a single post request.
</p>

## Getting Started
**Go must be installed!**
- Download the application from releases
- Run `go run .`, the server starts at port 9000

## Supported SQL
The database supports all the basic CRUD operations Examples on the supported sql syntax below.

### Create

<p align="justify">
    Table can be created with different attributes. Attributes must be inside parentheses. However there is no support for marking a primary key or to force any restrictions such as not null. The sql processor also supports only number and text values. <strong>Different attributes must be separated by a comma</strong>!
</p>

```sql
CREATE TABLE <table> (
    <column> INT,
    <column> VARCHAR(<length>)
)
```

### Read

<p align="justify">
    A single column or many columns can be requested from a table at the same time. items can be filtered with a where keyword and ordered with order by.
</p>

```sql
SELECT <column|*> FROM <table>
WHERE <condition>
ORDER BY <column> <ASC|DESC>
```

### Update

<p align="justify">
    A table items can be updated according to the filter in the where. A single item or many items can be updated at the same time. Attributes to be updated and their new values must be inside parentheses after set. <strong>Attributes mut be separated by a comma</strong>!
</p>

```sql
UPDATE <table>
SET (
    <name> = 10,
    <name> = 'Text'
)
WHERE <condition>
```

### Delete

<p align="justify">
    A single item or many items can be deleted with delete. Item filter can be specified with a where keyword.
</p>

```sql
DELETE <table>
WHERE <condition>
```

### Limitations
- `SELECT * FROM <table>` supports only selecting columns from a single table, many columns can be requested.
- `WHERE <condition>` supports only one condition, multiple where statements can however be combined. But testing value in range is possible, for example `0 <= x <= 1`. 
