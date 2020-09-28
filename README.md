Counter 
===

[![Build Status](https://img.shields.io/drone/build/0x76/counter?server=https%3A%2F%2Fdrone.xirion.net&style=for-the-badge)](https://drone.xirion.net/0x76/counter)
[![Docker Build](https://img.shields.io/docker/cloud/build/0x76/counter?style=for-the-badge)](https://hub.docker.com/r/0x76/counter)

Simple Go app for keeping count.

## Usage
The basic idea is that each path represents a key or counter which can be interacted with in a RESTful way.

Here is an example:

```sh
# First create the counter
curl -X POST localhost:8080/some/path
> { "/some/path": 0, "AccessKey": "acf38625-bc7b-4241-97be-55d4f20219f6" }

# Now we can increment it using the returned access key 
#   (the header of the response will also contain the key)
curl -X PUT -H "Authorization: Bearer acf38625-bc7b-4241-97be-55d4f20219f6" localhost:8080/some/path
> { "/some/path": 1 }

# Query it using GET
curl -X GET localhost:8080/some/path
> { "/some/path": 1 }

# We can also delete a counter if we want
curl -X DELETE -H "Authorization: Bearer acf38625-bc7b-4241-97be-55d4f20219f6" localhost:8080/some/path

# Now if we query it we'll get a 404 back
curl -X GET localhost:8080/some/path
> Counter not yet created
```
