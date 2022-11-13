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

### Router
* Server can return a router for route registering and run http server inside
```go
type Server interface {
	// Route return server root route
	Route() Router
	// Run attaches the router to a http.Server and starts listening and serving HTTP requests.
	// Note: this method will block the calling goroutine indefinitely unless an error happens.
	Run(addr ...string)
	// Close the server close, after this method called all request all failed with 500 status code
	Close()
	// Wait all server request handle event finish or timeout
	Wait(timeout int)
}
```
* Router is an interface you can handle http call response
. Router has a method SetRecover to set default recover actions
recover type is defined below:
```go
type Recover func(recV interface{}, req *http.Request, args []interface{}) (resp interface{}, err error)

```
Plate has already implemented router and server interface with gin.
In the ginroute package, we also define some struct tag control logic to handle request and response

*Request:*
* Format '{{tag_name}},{{loc}}:{{config option}}'
* Default Tag name is 'plate'
* loc include "header","body","path","form","query","mid". default location is body.
* config option has inline and full
  * inline will expand child object's fields  to parent fields. it's useful when parent will bring useless prefix path when marshal
  * full is valid only used in middle passing intermediate result. it means we take the current object as a full object and ignore child fields. it is very useful when some child fields is not public or can't be copied eg. time.Time  
  
*Response*
* response is default in json body with json tag
* support user define code and header to control http response code and headers
* with tag 'plate' loc 'resp' field is code and header. code type must be number type and headers must be map[string]string. user header empty value to delete a header key.
* eg.
```
  type Resp struct {
	Data struct {
		Items []*model.Data    `json:"data"`
		Total int              `json:"total"`
	} `json:"data"`
	Code   int               `json:"-" plate:"code,resp"`
	Header map[string]string `json:"-" plate:"header,resp"`
  }
```

## How To Use
You can find some way from test cases and examples right now.