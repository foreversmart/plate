### plate
plate is a micro server framework. which provide simple micro service organization and replaceable components. 

* code as config so all the config is explicitly defined by code
eg.

```cassandraql
    json type config
    
    {
      "host": "127.0.0.1",
      "port": "9099"
    }
    
    defined code as
    type ServerConfig struct {
      Host string `json:"host"`
      Port string `jsong:"port"`
    }

```



* code as logic

##### Design

* Config Component : as we say config is explicitly defined by code but how?

* Configer interface represent an interface of a config loader. Through 
Configer we can easy defined what kind of config file or config center,
conveniently register concrete config struct to it

* Config structure: if a component which composed of one or several packages wants
config something. we suggest the component hold a config struct which
defined the specific config content only need by it. then then caller will use 
configer to init the config struct content data.

```
type Configer interface {
	Init(mode ModeType, path, configName, host, meta string) error
	Register(key string, config interface{})
}
``` 
