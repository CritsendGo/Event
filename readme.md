# MODULE EVENT
Event module is used on different subproject to normalise and create a secure buffer

## Installation



```bash
go get github.com/CritsendGo/modEvent
```
Proxy via Nginx to manage https

## Usage
```go
import("github.com/CritsendGo/modEvent")

func main(){
    event:=&modEvent.Event{
        Code : 100,
        Detail : "Example",		
    }
    modEvent.AddEvent(event)
}

```
## License
Attribution-NonCommercial-NoDerivatives 