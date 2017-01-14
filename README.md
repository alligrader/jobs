# jobs
This repository contains helpful functions and types for specifying autograder jobs.

Reminder to self: You can use `bash -c ` to execute a command as the third argument.

So the go code:

```
   cmd := exec.Command("bash", "-c", commandString)
    if err := cmd.Start(); err != nil {
        log.Fatal(err)
    }
    done := make(chan error)
    go func() { done <- cmd.Wait() }()
    select {
    case err := <-done:
        // exited
    case <-time.After(10*time.Second):
        // timed out
    }
```

runs the command string and waits 10 seconds for it to time out, where command string is in quotes (if you pass it as a single argument, it might be implicitly in quotes. Not 100% on that one.)
