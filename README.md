# pivot-lang-coderunner
a pivot lang coderunner written with go

## usage
```go
var req struct {
		Code string `json:"code"`
	}
```
As the code above shows, the code you want to run should be in the `code` field of the request body as a string. The method should be `POST` and the content type should be `application/json`.
```go
type Response struct {
	Status        int    `json:"status"`
	CompileOutput string `json:"compileOutput"`
	RunOutput     string `json:"runOutput"`
}
```
The response is a json object with three fields: `status`, `compileOutput` and `runOutput`. The `status` field is the status code of the coderunning. It is 1 if compiled correctly and 0 if not. The `compileOutput` field is the output of the compiler. The `runOutput` field is the output of the execution of the code.

## features

### update.go
The container image is updated regularly. Older images will be deleted.

### cache.go
Code and its response will be inserted into cache. Later requests with the same code will get the response from the cache.

### runner.go
The code is compiled and executed in a container. The container is created and deleted every time the code is run.

### initAndExit.go
Init: create a container and relevant file and directory. 
Exit: delete all containers when the program exits.

### sync
The server can handle multiple requests at the same time. However, if the requests exceed a certain number, it shall wait for the previous requests to finish. 

## variables
```go
const (
	defaultContainerName = "coderunner"
	cacheTTS             = 10
	containerNum         = 10
)
```
`defaultContainerName` is the name of the container. 

`cacheTTS` is the time to live of the cache. 

`containerNum` is the number of containers that can be created at the same time.

```go
var (
	containerIndex = uint64(0)
	imageName      = "registry.cn-hangzhou.aliyuncs.com/pivot_studio/pivot_lang"
	imageTag       = "latest"
	updateInterval = 7 * 24 * time.Hour
)
```
`containerIndex` is the index of the container starting from 0.

`imageName` is the name of the image.

`imageTag` is the tag of the image.

`updateInterval` is the interval of the update of the image.