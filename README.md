# Database SQL
![Version](https://img.shields.io/badge/go-1.21+-blue.svg?style=flat)
![License](https://img.shields.io/badge/license-MIT-green.svg?style=flat)

<p align="justify">
    An extremely simple local sql database that supports only a few of the sql keywords. The data can be requested/updated from the server by http post request at root. The data is send in json format as http response if the request sql is in correct format. The application supports only a single sql query send a single post request.
</p>

**NOT READY TO BE USED YET!!!**

## Getting Started
Clone the repository.

```
git clone https://github.com/WilliwadelmaWisky/DatabaseSQL.git
```

cd to source directory.

```
cd DatabaseSQL/src
```

Start the program. 

```
go run . [URI] [PORT]
```

> [!NOTE]
> URI can be specified (defaults to **default**). URI defines to location of the database files on the disk, nested directories should be separeated by '/' on all operating systems.
> Port can be specified (defaults to **9000**). Port is checked from the second argument so URI must be specified before.

Communicate with the database via curl or some other tool.

```
curl -i -X POST -d "SELECT * FROM table" localhost:9000
```

> [!NOTE]
> If you specified a port use that instead of *9000*.

## Features
<p align="justify">
    The database supports all the basic CRUD operations Examples on the supported sql syntax below. If you want to use values with spaces, only single quotes are supported. Syntax allows additional spaces/newlines in the sql expressions and is not case sensitive.
</p>

### Create a new Table

<p align="justify">
    Table can be created with different attributes. Attributes must be inside parentheses. However there is no support for marking a primary key or to force any restrictions such as not null. The sql processor also supports only number and text values. Let's create a table named <i>artists</i> as an example. A single artist in a table contains an <i>id</i> (int), a <i>name</i> (string) and an <i>age</i> (int).
</p>

```sql
-- Create a table of which items have id, name and age
CREATE TABLE artists (
    id INT,
    name VARCHAR,
    age INT
)
```

> [!NOTE]
> Different attributes must be separated by a comma and must be specified inside parentheses!

### Insert data to a table

<p align="justify">
    Data can be inserted to a table with basic sql insert into syntax. Every attribute does not have to be explicitly typed. If an attribute is not inserted it will be initialized to a default value. Let's insert some data to the <i>artists</i> table we created above.
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

> [!NOTE]
> Only single sql expression can be sent at a time so doing the previous four expression cannot be sent at the same time. Different attributes must be separated by a comma and parentheses must be used! Note that last artist will now have default age of 0.

### Fetch data from a table

<p align="justify">
    Data can be fetch from a single table with basic sql select syntax. A single column or many columns can be requested from a table at the same time. items can be filtered with a where keyword and ordered with order by. Let's fetch some data from the <i>artists</i> table created above.
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

> [!NOTE]
> Select supports only selecting columns from a single table, however many columns can be requested separated with comma. Where supports only one condition, multiple where statements can however be combined. But testing value in range is possible, for example `0 <= x <= 1`. 

For the first select expression returned data is in the following format.

```json
{
    "column_names": ["id", "name", "age"],
    "rows": [
        ["1", "Artist 1", "50"],
        ["2", "Artist 2", "25"],
        ["3", "Artist 3", "45"],
        ["4", "Artist 4", "0"],
    ],
    "row_count": 4
}
```

### Update data in a table

<p align="justify">
    A data can be updated with update syntax. A filter, to which items are updated, can be specified with where syntax. A single item or many items can be updated at the same time. Attributes to be updated and their new values must be inside parentheses after set. Let's update the <i>artists</i> table as an example.
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

> [!NOTE]
> Attributes must be separated by a comma and parentheses are important!

### Delete data from a table

<p align="justify">
    A single item or many items can be deleted with delete. Item filter can be specified with a where keyword. Let's delete some data from the <i>artists</i> table created above.
</p>

```sql
-- Delete all
DELETE FROM artists

-- Delete artist that has id=1
DELETE FROM artists
WHERE id = 1
```
