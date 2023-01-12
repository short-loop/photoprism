Steps to use this sdk locally in your go service

1. In go.mod file use replace to replace the package names of sdk and common with your local path  
   (this step will be required only when using local sdk)
    ```
    replace github.com/short-loop/shortloop-sdk-go => D:\shortloop\shortloop-sdk-go
    replace github.com/short-loop/shortloop-common-go => D:\shortloop\shortloop-common-go
    ```

2. Include the sdk in your go.mod file
    ```
    require github.com/short-loop/shortloop-sdk-go v0.0.1
    ```

3. Run go mod tidy
    ```
    go mod tidy
    ```

4. Start using the sdk in your code  

    For gin:
    ```
    import "github.com/short-loop/shortloop-sdk-go/shortloopgin"
    ```
   
    For mux:
    ```
    import "github.com/short-loop/shortloop-sdk-go/shortloopmux"
    ```

5. Initialize the sdk  
   Example for gin:
    ```
    router := gin.Default()
    sdk, err := shortloopgin.Init(shortloopgin.Options{
        ShortloopEndpoint: "http://localhost:8080",
        ApplicationName:   "test-service-go",
        LoggingEnabled:    true,
        LogLevel:          "INFO",
    })
    if err != nil {
        fmt.Println("Error initializing shortloopgin: ", err)
    } else {
        router.Use(sdk.Filter())
    }
    ```
   Example for mux:
    ```
    mux := mux.NewRouter()
    sdk, err := shortloopmux.Init(shortloopmux.Options{
        ShortloopEndpoint: "http://localhost:8080",
        ApplicationName:   "test-service-go",
        LoggingEnabled:    true,
        LogLevel:          "INFO",
    })
    if err != nil {
        fmt.Println("Error initializing shortloopmux: ", err)
    } else {
        mux.Use(sdk.Filter)
    }
    ```

