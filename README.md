tcpchan
=======

A network communication channel for Google Go

## Usage
Server:

    ch, err := tcpchan.Listen("<local ip>:8976")
    if err != nil {
        ...
    }
    
    ch <- "Test"
    fmt.Println(<-ch)
    close(ch)
    
Client:

    ch, err := tcpchan.Dial("<remote ip>:8976")
    if err != nil {
        ...
    }
    
    fmt.Println(<-ch)
    ch <- "Hello!"
    
## Known issues
When the application is terminated before all data is sent the remaining data is lost.

## License
Zlib-License
