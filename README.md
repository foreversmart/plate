# Plate
plate is a micro server framework. which provide simple micro service organization and replaceable components. 


### Design
Plate wants to make every piece of the framework easy to plug and convenient for use. 
Just like the framework's name plate we can place all the plate component together and also 
convenient contains logic.

### Philosophy
* All the plate component plug with Go interface, 
And every component provide a default way to use which is out of the box.
which means you can use component without customize or initialize it.

* Code as logic all the logic must be explicitly defined by code
include config. This means there is no logic hide at code level, 
What you can see from code is just the code build and run. It will 
improve code readability and easy to debug your code.

### GO TO DETAIL
* A Plate Component can be everything
* Component has a frontend interface for caller use
* Component has at least one default implementation and a default way to use just like Go http DefaultClient
* Customize a component with a standard Config Interface, 
Since most of the components is different and want some key/value to config
```go
type Configer interface {
	Init(mode ModeType, path, configName, host, meta string) error
	Register(key string, config interface{})
}
``` 

## Component Design
### Config 
* Config Component : as we say config is explicitly defined by code but how?

* Configer interface represent an interface of a config loader. Through
  Configer we can easy defined what kind of config file or config center,
  conveniently register concrete config struct to it

* Config structure: if a component which composed of one or several packages wants
  config something. we suggest the component hold a config struct which
  defined the specific config content only need by it. then then caller will use
  configer to init the config struct content data.

```go
type Configer interface {
	Init(mode ModeType, path, configName, host, meta string) error
	Register(key string, config interface{})
}
``` 


## How To Use
You can find some way from test cases and examples right now.