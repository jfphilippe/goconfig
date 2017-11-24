# Package goconfig
implements a configuration reader that can read valeurs form :
- json files
- txt files
- env variables

values are stored as "nested" maps.

## Quick Usage

```go
builder := NewBuilder("CTX_", nil)
// Main config values
_, err := builder.LoadJsonFile("/etc/myapp/config.json")

// complete config with values from another file
// existing values will NOT be updated.
_,err := builder.LoadTxtFile("/etc/default/mayapp.txt")

// Use it
config := builder.GetConfig()

name := config.GetString("name", "default name")
num_proc := config.GetInt("proc.limit", 5)

dbconfig := config.GetConfig("database")

db_url := dbconfig.GetString("url", "")
db_port := dbconfig.GetInt("port", 1234)
...
```
## Value Expansion
GoConfig can expand values identified by ${ }.

Exemple :
```
db.user = john

database.url=${db.ser}@/dbname
```

In the previous configuration, _database.url_ will be resolved as : john@/dbname


## Text File Format
It's a basic file format where each value is writen in a line.

### Comments
Comments are writen  with a # char at begining of the line (first non white space char). Comments end at the end of the line.

Exemple :

```
# this is a Comments
   # this is another one
```

### values

Syntax :
```
key = value
```

Where _value_ is any string that ends at the end of the current line. Multi lines values are not supported.
_key_ is the "name" of the value. Nested names are separated with a '.', i.e. : database.name
