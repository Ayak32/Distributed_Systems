1. What is the difference between an unbuffered channel and a buffered channel?
    - A buffered channel specifies a capasity for the channel while an unbuffered channel has no capacity.

2. Which is the default in Go? Unbuffered or buffered?
    - unbuffered

3. What does the following code do?
    - This code:
        1) Creates a channel 'ch' that recieves strings sends the string "hello world!" into the channel
        2) Recieves the string "hello world!" and stores it in the variable message
        3) Prints the string stored in message

4. In the function signature of MergeChannels in merge_channels.go, what is the difference between <-chan T, chan<- T and chan T?
    - chan T is bidirectional (you can send and recieve from the channel)
    - chan<- T allows you to only recieve from the channel
    - <-chan T allows you to only send from the channel

5. What happens when you read from a closed channel? What about a nil channel?
    - 'ok' will be set to false if there are no more values to receive and the channel is closed.
    - reading from a nil channel will block until something is sent into the channel and it is ready to be recieved

6. When does the following loop terminate?
    - the loop will terminate when the channel is closed

7. How can you determine if a context.Context is done or canceled?
    - we can use the Done() method, which returns a channel that is closed when the context is done or canceled

8. What does the following code (most likely) print in the most recent versions of Go (e.g., Go 1.23)? Why is that?
    - it would most likely print "all done! 4 4 4" because
    1) go routines are non-blocking, so  "all done!" is printed right after launching the go routines and before the goroutines finish sleeping.
    2) The goroutines capture i by reference and not by value, so they will print out the same value
    3) The loop finishes very quickly, so by the time the goroutines are executed after sleeping, I is already 4

9. What concurrency utility might you use to "fix" question 8?
    - using a WaitGroup would allow all the goroutines to finishing executing before the main function prints "all done!"

10. What is the difference between a mutex (as in sync.Mutex) and a semaphore (as in semaphore.Weighted)?
    - a mutex allows only one goroutine at a time, while a semaphore can allow multiple goroutines to access a resource


11. What does the following code print?
    - nil
    - 0
    - true
    - ""
    - 0
    - nil
    - {}


12. What does struct{} in the type chan struct{} mean? Why might you use it?
    - it means the channel that sends and receives empty structs
    - since the struct takes zero bytes, the channel can be used as a synchronization mechanism without needing to pass any data

