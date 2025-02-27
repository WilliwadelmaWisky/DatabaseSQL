# Database SQL
![Version](https://img.shields.io/badge/go-1.21+-blue.svg?style=flat)
![License](https://img.shields.io/badge/license-MIT-green.svg?style=flat)

<p align="justify">
    An extremely simple local sql database that supports only a few of the sql keywords. The data can be requested and updated on the server by http post request at root. The data is send back in json format as http response if the request sql is in allowed format. The application supports only a single sql query send in a single post request. The database server has cors enabled so it can be used together with a javascript frontend application for example.
</p>

> [!WARNING]
> There are still some errors when sql syntax is invalid, but the database is usable.

## Getting Started
1. Clone the repository.
2. Run the start script or optionally run the application manually with `go run`.
3. Communicate with the database via curl or some other tool. The server has CORS enabled so the database can also be accessed through a web application for example.

```
git clone https://github.com/WilliwadelmaWisky/DatabaseSQL.git
```

```
./DatabaseSQL/DatabaseSQL.sh <DIRECTORY> [PORT]
```

```
curl -X POST -d "SELECT * FROM table" localhost:9000
```

> [!IMPORTANT]
> DIRECTORY is mandatory so it must be excplicitly specified. The DIRECTORY defines the location of the database files on the disk, nested directories should be separeated by '/' on all operating systems. The absolute path on the is`/home/user/DIRECTORY/` for linux and `C:\Users\user\DIRECTORY\` for windows. The DIRECTORY is checked from the first argument.
> PORT is optional so it can be specified (defaults to **9000**). The server is activated on `localhost:PORT`. The PORT is checked from the second argument.
> More information can be accessed with `[-h | --help] [-v | --version]` flags as the first argument.

## Features
<p align="justify">
    The database supports all the basic CRUD operations Examples on the supported sql syntax below. Syntax allows additional spaces/newlines in the sql expressions and is not case sensitive. The database http server has cors enabled so browser can communicate with the server.
</p>

> [!IMPORTANT]
> If you want to use values with spaces, only single quotes are supported `'value with spaces'`.

### Create a new Table
<p align="justify">
    Table can be created with different attributes. Attributes must be inside parentheses. However there is no support for marking a primary key or to force any restrictions such as not null. The sql processor also supports only number and text values. Let's create a table named <i>artists</i> as an example. A single artist in a table contains an <i>id</i> (int), a <i>name</i> (string) and an <i>age</i> (int). Data is saved automatically on disk after the table gets created.
</p>

```sql
-- Create a table of which items have id, name and age
CREATE TABLE artists (
    id INT,
    name VARCHAR,
    age INT
)
```

> [!IMPORTANT]
> Different attributes must be separated by a comma and must be specified inside parentheses! Table names must be unique.

### Insert data to a table
<p align="justify">
    Data can be inserted to a table with basic sql insert into syntax. Every attribute does not have to be explicitly typed. If an attribute is not inserted it will be initialized to a default value. Let's insert some data to the <i>artists</i> table we created above. Data is saved automatically on disk after data is inserted.
</p>

```sql
-- Insert Artist 1
INSERT INTO artists (id, name, age) 
VALUES (1, 'Artist 1', 50)

-- Insert Artist 2
INSERT INTO artists (id, name, age) 
VALUES (2, 'Artist 2', 25)

-- Insert Artist 3
INSERT INTO artists (id, name, age) 
VALUES (3, 'Artist 3', 45)

-- Insert Artist 4, default age
INSERT INTO artists (id, name) 
VALUES (4, 'Artist 4')
```

> [!IMPORTANT]
> Only single sql expression can be sent at a time so doing the previous four expression cannot be sent at the same time. Different attributes must be separated by a comma and parentheses must be used! Note that last artist will now have default age of 0.

### Fetch data from a table
<p align="justify">
    Data can be fetch from a single table with basic sql select syntax. A single column or many columns can be requested from a table at the same time. items can be filtered with a where keyword and ordered with order by. Let's fetch some data from the <i>artists</i> table created above. Does not load or save data on disk.
</p>

```sql
-- Get the whole table
SELECT * FROM artists

-- Get all artists of age>40
SELECT * FROM artists
WHERE age > 40

-- Get artists name and age that are in their forties, also order by name ascending
SELECT name, age FROM artists
WHERE 40 <= age <= 49
ORDER BY name ASC

-- Get artists name that are age<50, also order by name ascending and if same name then order by age descending
SELECT name FROM artists
WHERE age < 50
ORDER BY name ASC
ORDER BY age DESC
```

> [!IMPORTANT]
> Select supports only selecting columns from a single table, however many columns can be requested separated with comma. Where supports only one condition, multiple where statements can however be combined. But testing value in range is possible, for example `40 <= x <= 49`. When comparing to single value, for example `age > 40`, table name must be on the left side of the operator.

For the first select expression returned data is in the following format.

```json
{
    "columns": ["id", "name", "age"],
    "column_types": ["INT", "VARCHAR", "INT"],
    "data": [
        ["1", "Artist 1", "50"],
        ["2", "Artist 2", "25"],
        ["3", "Artist 3", "45"],
        ["4", "Artist 4", "0"],
    ]
}
```

### Update data in a table
<p align="justify">
    A data can be updated with update syntax. A filter, to which items are updated, can be specified with where syntax. A single item or many items can be updated at the same time. Attributes to be updated and their new values must be inside parentheses after set. Let's update the <i>artists</i> table as an example. Data is saved automatically on disk after data is updated.
</p>

```sql
-- Update ages of every artist to 50
UPDATE artists (age)
VALUES (50)

-- Update artist that has id=1, age to 60
UPDATE artists (age)
VALUES (60)
WHERE id = 1

-- Update artist that has id=1, name to Artist 11 and age to 60
UPDATE artists (name, age)
VALUES ('Artist 11', 60)
WHERE id = 1
```

> [!IMPORTANT]
> Attributes must be separated by a comma and parentheses are important!

### Delete data from a table
<p align="justify">
    A single item or many items can be deleted with delete. Item filter can be specified with a where keyword. Let's delete some data from the <i>artists</i> table created above. Data is saved automatically on disk after data is deleted.
</p>

```sql
-- Delete all
DELETE FROM artists

-- Delete artist that has id=1
DELETE FROM artists
WHERE id = 1
```

### Delete an existing table
<p align="justify">
    A table can be deleted from the database. Let's delete the <i>artists</i> table created above. Data is saved automatically on disk after table is deleted.
</p>

```sql
-- Delete the artists table
DROP TABLE artists
```

### Getting database metadata
<p align="justify">
    The database metadata can be requested from `localhost:9000/information_schema` as a http get request (change the <i>9000</i> to correct port). 
</p>

The metadata is send in a following format as json.

```json
{
    "tables": ["artists"]
}
```
