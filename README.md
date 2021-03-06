# Package goconfig

implements a configuration reader that can read values form :

- json files
- txt files
- env variables

values are stored as "nested" maps.

## Quick Usage

```go
// Create a Builder, env vars are searched in uppercase with CTX_ prefix.
builder := NewBuilder("CTX_", nil)
// Main config values
_, err := builder.LoadJsonFile("/etc/myapp/config.json")

// complete config with values from another file
// existing values will NOT be updated.
_,err := builder.LoadTxtFile("/etc/default/mayapp.txt")

// Or
// _,err := builder.LoadFiles("/etc/myapp/config.json","/etc/default/mayapp.txt")

// Use it
config := builder.Config()

name := config.GetString("name", "default name")
num_proc := config.GetInt("proc.limit", 5)

dbconfig := config.GetConfig("database")

db_url := dbconfig.GetString("url", "")
db_port := dbconfig.GetInt("port", 1234)

// or
db_url = config.GetString("database.url","")
db_port = config.GetInt("database.port", 1234)
...
```

## Value Expansion

### Basic Expension

GoConfig can expand values identified by `${ }`.

Exemple :

```txt
db.user = john

database.url=${db.ser}@/dbname
```

In the previous configuration, _database.url_ will be resolved as : __john@/dbname__

### Nested Expension

Expensions can be "nested".

Exemple :

```txt
dev.db.pwd=azerty
int.db.pwd=qwerty

...
env = dev
...

database.pwd = ${${env}.db.pwd}
```

Here _database.pwd_ will be resolved as __azerty__.

### Recursive expension

Expension may be resolved recursively

Exemple :

```txt
file.root=/tmp/myapp
file.conf.txt=${file.root}/conf.txt

config.file=${file.conf.txt}
```

_config.file_ should be translated into __/tmp/myapp/conf.txt__

## Text File Format

It's a basic file format where each value is writen in a line.

### Comments

Comments are writen  with a `#` char at begining of the line (first non white space char). Comments ends at the end of the line.

Exemple :

```txt
# this is a Comments
   # this is another one
```

### values

Syntax :

```txt
key = value
```

Where _value_ is any string that ends at the end of the current line. Multi lines values are not supported.
_key_ is the "name" of the value. Nested names are separated with a `.` (dot), i.e. : `database.name`
