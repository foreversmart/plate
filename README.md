# plate
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
