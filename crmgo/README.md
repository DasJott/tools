# CleverReach Mongo Adapter
Provides a ready to go mgo instance with very few lines.<br>
Makes sure your connections never dies.

## Environment / crconfig variables
- MONGO_HOST (default "localhost")
- MONGO_PORT (default "27017")

### Multiple connections (example)
- MONGO_HOST_SPECIAL
- MONGO_PORT_SPECIAL

## Usage
```go

type FamilyMember struct {
    Name string `bson:"name"`
    Age  int    `bson:"age"`
}

// Do adds Mary and George to 'family' collection in database 'my_db'
func Do() {
    db := crmgo.MustOpen("my_db")
    defer db.Close()

    db.C("family").Insert(&FamilyMember{
        Name: "Mary",
        Age: 32,
    })

    db.C("family").Insert(&FamilyMember{
        Name: "George",
        Age: 33,
    })
}
```
As you can see the `C()` function returns the according `*mgo.Collection` object to work on.<br>

**Note:** Make sure you never store the collection, but rather make a wrapper for getting it:
```go
func getFamily() *mgo.Collection {
    return db.C("family")
}
```
Otherwise you are not using the code that checks the connection avaibility!

### Using multiple mongo connections
If you need multiple mongo connections, you can also configure them just fine using suffixes.<br>
The following example shows the configuration of two seperate connections.

The config file would look like this:
```ini
MONGO_HOST_FIRST = my.firsthost.com
MONGO_PORT_FIRST = 27017

MONGO_HOST_SECOND = my.secondhost.com
MONGO_PORT_SECOND = 27017
```

And the code like this:
```go
func Do() {
    db1 := crmgo.WithSuffix("first").MustOpen("my_db")
    defer db1.Close()

    db2 := crmgo.WithSuffix("second").MustOpen("my_db")
    defer db2.Close()

    // ...
}

```

## Query helper Q
There are some functions you can use to make life easier on writing queries.
Q is for 'query' and can be used in `Find()` on a collection.<br>
Q has some functions you can use:

- **GT(key string, val interface{}) Q**<br>
  GT is greater than

- **GTE(key string, val interface{}) Q**<br>
  GTE is greater than or equal

- **LT(key string, val interface{}) Q**<br>
  LT is less than

- **LTE(key string, val interface{}) Q**<br>
  LTE is less than or equal

You can also use same functions to actually retrieve Q:

- **GT(val interface{}) Q**<br>
  GT is greater than

- **GTE(val interface{}) Q**<br>
  GTE is greater than or equal

- **LT(val interface{}) Q**<br>
  LT is less than

- **LTE(val interface{}) Q**<br>
  LTE is less than or equal

As every of those functions return a Q they are chainable to add multiple conditions.

## Logger
You can set a logger to get debug information.<br>
The logger must implement the DBLogger interface:
```go
	DBLogger interface {
		Debugln(...interface{})
	}
```
The default logger simply uses `fmt.Println`.