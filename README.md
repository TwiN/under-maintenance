# under-maintenance

[![Build Status](https://travis-ci.com/TwinProduction/under-maintenance.svg?branch=master)](https://travis-ci.com/TwinProduction/under-maintenance)
[![Coverage Status](https://coveralls.io/repos/github/TwinProduction/under-maintenance/badge.svg?branch=master)](https://coveralls.io/github/TwinProduction/under-maintenance?branch=master)
[![Docker pulls](https://img.shields.io/docker/pulls/twinproduction/under-maintenance.svg)](https://cloud.docker.com/repository/docker/twinproduction/under-maintenance)

A very small Docker image that returns `Under maintenance` for every request made on the port 80 by default. 
The content returned can be customized, see [Page Content](#page-content).

By default, the status code returned will by 503 (Service Unavailable), but it can be customized through environment variable.


## Table of Contents

- [under-maintenance](#under-maintenance)
  * [Table of Contents](#table-of-contents)
  * [Usage](#usage)
  * [Why?](#why)
  * [Customization](#customization)
    + [Page content](#page-content)
    + [Status code](#status-code)
    + [Retry-After](#retry-after)
  * [Specifications](#specifications)


## Usage

Pull the image from Docker Hub:

```
docker pull twinproduction/under-maintenance:latest
```

Run it:

```
docker run --name under-maintenance -p 0.0.0.0:80:80 twinproduction/under-maintenance
```


## Why?

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

That being said, this is just one of many ways to implement it.


## Customization

### Page content

If you're looking to return more than just a plain `Under maintenance`, you can. All you have to do is overwrite the `under-maintenance.html` file situated in the same directory where the binary is executed, whether it be by using `configs` or `volumes`, like so:

```yaml
version: '3.7'
services:
  fallback:
    image: twinproduction/under-maintenance
    restart: always
    ports:
      - 80:80
    volumes:
      - ./under-maintenance.html:/under-maintenance.html
```

**NOTE**: To reduce the footprint of the application, the content of the `under-maintenance.html` file is not read every time a request is received. 
Instead, it is read only once when the application starts.


### Status code

You can modify the status code by setting the `UNDER_MAINTENANCE_STATUS_CODE` environment variable to a status code of your choice.

By default, the status code is 503, and truthfully, it should remain so.

According to [RFC 2616 section 10](https://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html):

> **10.5.4 503 Service Unavailable**
> The server is currently unable to handle the request due to a temporary overloading or maintenance of the server. The implication is that this is a temporary condition which will be alleviated after some delay. If known, the length of the delay MAY be indicated in a Retry-After header. If no Retry-After is given, the client SHOULD handle the response as it would for a 500 response. 

Returning 503 is supposedly also important if you're planning on having an extended maintenance period, otherwise, it could affect your search ranking.


### Retry-After

You can modify the value of the `Retry-After` header by setting the `UNDER_MAINTENANCE_RETRY_AFTER` environment variable. 

The `Retry-After` header indicates how long the service is expected to be unavailable.
This is mostly important for crawling bots.

There are two main use cases:

1. **The application is being redeployed or small changes are being made.** 
In this case, you should use a fixed number. The unit is seconds and by default, the value is `300`. This will indicate that the website will be back up momentarily.

2. **The application is ongoing massive changes, and will be down for longer than a day.**
In this case, you should use a date, more precisely, the date at which the maintenance is planned to be over.

For more information, see https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Retry-After

**NOTE**: The `Retry-After` header will not be added if the status code is not 503 (default) or 429.


## Specifications

| Property    | Value |
|------------:|-------|
| language    | Go    |
| port        | 80    |
| image size  | 6.5MB |