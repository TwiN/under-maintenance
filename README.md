# under-maintenance

A very small Docker image that returns `Under maintenance` for every request made on the port 80.


### Why?

Let's say your application goes down and you have no replicas. 
Rather than leaving your visitors in the dark, you can use this to let them know that 
your website is _"under maintenance"_

As for how to apply it, there are many ways. The more obvious one is having a reverse proxy point `under-maintenance` to the same route as your main application, but with a lower priority. That way, 
`under-maintenance` will take over when your application goes down.

On the flowchart below, all traffic is going to `main application`, because its priority is higher
than the priority assigned to `under-maintenance`.

```
                                    +-----------------+
                             P100   |                 |
                           +------> |main application |
           +---------------+        |                 |
 traffic   |               |        +-----------------+
+------->  | reverse proxy |
           |               |        +-----------------+
           +---------------+        |                 |
                           +------> |under-maintenance|
                              P1    |                 |
                                    +-----------------+
```

On the flowchart below, all traffic is going to `under-maintenance`, because the `main application` is not deployed.

```
           +---------------+
 traffic   |               |
+------->  | reverse proxy |
           |               |        +-----------------+
           +---------------+        |                 |
                           +------> |under-maintenance|
                              P1    |                 |
                                    +-----------------+
```

On the flowchart below, all traffic is going to `under-maintenance`, because the application is unhealthy. Obviously, you have to implement health checks to make this work.

```
                                    XXXXXXXXXXXXXXXXXXX
                             P100   XXXXXXXXXXXXXXXXXXX
                           +------> XXXXXXXXXXXXXXXXXXX
           +---------------+        XXXXXXXXXXXXXXXXXXX
 traffic   |               |        XXXXXXXXXXXXXXXXXXX
+------->  | reverse proxy |
           |               |        +-----------------+
           +---------------+        |                 |
                           +------> |under-maintenance|
                              P1    |                 |
                                    +-----------------+
```


### Specifications


| Property   | Value |
|-----------:|-------|
| language   | Go    |
| port       | 80    |
| image size | ~8MB  |
