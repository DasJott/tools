# CleverReach Config
This is a very simple config reader.<br>
It reads the config from cli switch, environment, from config file or uses a given default.<br>
The format of the config file is same as for environments.

## Usage
The usage from within your code is rather easy:

```ini
# this is your configfile e.g. config.env
    MY_URL = https://my.productive.url
    MAGIC_NUM = 42
    MAGIC_FLOAT = 42.23
    SIMPLE = true
```

```go
func myConfig() {
    // Read a config file if you wish
    // error indicates mostly 'file not existing'
    err := crconfig.Read("config.env")
	if err != nil {
		fmt.Println(err.Error())
	}

    // Get simply gets the according string
    myURL := crconfig.Get("MY_URL", "https://find.me.here")

    // GetInt gets an int
    meaning := crconfig.GetInt("MAGIC_NUM", 42)

    // GetFloat gets a float
    no_meaning := crconfig.GetFloat("MAGIC_FLOAT", 42.56)

    // GetBool - surprise - gets a bool
    simple := crconfig.GetBool("SIMPLE", true)
}
```
### Bind config to your struct
You can use `Bind()` as often as you wish, e.g. to get small portions af the config in different packages.<br>

Example:
```go
// Config takes the config
// it defines which values to use and also provides default values by tag
type Config struct {
    MyURL string `env:"MY_URL,https://my.productive.url"`
    MagicNum = 42 `env:"MAGIC_NUM,42"`
    MagicFloat = 42.23 `env:"MAGIC_FLOAT,42.23"`
    Simple = true `env:"SIMPLE,true"`
}

func (c *Command) Init() {
    // Bind fills a struct with values
    obj := MyStruct{}
    crconfig.Bind(&obj)

}
```
**You can also use `BindExclusive()` to simply ignore all fields of your struct, not having an `env` tag.**<br>
This way you can even use your `Command` for holding config.

## File format
You can call the file e.g. `config.env` or whatever you want, as long as you specify it correctly on `Read()`.<br>
A valid file looks like the following:
```ini
# this is my cool config file
MY_URL=https://find.me.there
MAGIC_NUM=42
SIMPLE=true
```
As you can see, the file uses the environment variable style.<br>
You can use spaces around the equal sign though, if you need it.<br>
Comments can be made starting the line with a `#` _at the very begining_.


## Command Line Switches
For starting your app with certain parameters on the fly, use switches.<br>
Switches are first choice, even before environments, config and default anyways.<br>
You can simply reference switches to config values as follows:
```go
MAGIC_NUM=42
SIMPLE=true

# switches refering to values
-n MAGIN_NUM
-s SIMPLE
```

## Advantages
- It automaticly reads the environment values if existing (e.g. running in Docker)
- You can have personal local settings within a config file.
- Default values used if nothing set in environment or in a file.
- Config file is a place to find used values for usage in yml files.
- Very lightweight with quick parsing and processing.
