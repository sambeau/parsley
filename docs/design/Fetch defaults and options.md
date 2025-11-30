# Fetch defaults and options

Parsely is mostly a server side web-rendering language. It is trying to make the basic stuff super-simple, the slightly complicated simple and the complicated stuff can be complicated. It's 'batteries included' as they say. As such, should alway default to the basic common web defaults.

so, for the fetch object,

- ``<=/=`` should default to ``GET``
- ``=/=>`` should default to ``POST``

As those are the de-facto standards of the web

## Direction of the arrows

The direction of the arrows should match the direction of the **significant** message: the data flow.

- ``<=/=`` Server sends data to client
- ``=/=>`` Client sends data to Server

The request isn't the significant data; the PAYLOAD is.

So,

- ``PUT`` sends a payload to the server which means ``=/=>``
- ``PATCH`` sends a payload to the sever which means ``=/=>``

### DEkete

For delete, the import message is the request for the file to be deleted not the response

So,

- DELETE sends an important message which means ``=/=>``

## Specifying HTTP Method

The vast majotity of HTTP calls are GET. Next is POST. Many APIs use POST for all CRUD/REST operations. However, we have to support *all* of the 5 HHTP methods, but we can be a little less intuitive about it. Firstly, we can add a dictionary of Fetch option. This is fine though a little ugly, plus there is a danager of a mismatch between the operatod and the options — e.g. `<=/=``with a ``POST``. The horror!

If possible, we should match (and limit) the HTTP method in the settings to the corresponding arrow operator.

## Shortcuts for PUT, PATCH and DELETE

We could supply shortcut values/methods on the fetch object to specify the HTTP method. This would match to what we do with SFTP.

SFTP can do this:

```parsley
{result, err} <=/= conn(@/data.json).json
```
We should let HTTP do something similar:

```parsley
let updated = bob =/=> JSON(@https://api.example.com/users/1).put
```

❌ This example from the refernce isn't what we need: the arrow is going the wrong way and the payload should be on the client side of the arrow:

```parsley
// POST with JSON body
let response <=/= JSON(@https://api.example.com/users, {
    method: "POST",
    body: {name: "Alice", email: "alice@example.com"},
}) 
```

✅ This is so much clearer

```parsley
let response = {name: "Alice", email: "alice@example.com"} =/=> JSON(@https://api.example.com/users) // defaults to POST
```

And if you want it to be a PUT or a PATCH

```parsley
let response = {name: "Alice", email: "alice@example.com"} =/=> JSON(@https://api.example.com/users).put
let response = {name: "Alice", email: "alice@example.com"} =/=> JSON(@https://api.example.com/users).patch
```

If need-be we could use a method call, e.g.

```parsley
let response = {name: "Alice", email: "alice@example.com"} =/=> JSON(@https://api.example.com/users).put()
let response = {name: "Alice", email: "alice@example.com"} =/=> JSON(@https://api.example.com/users).patch()
```

But, I'm less keen as it's not clear what the precedence is, and might confuse users — does the patch apply after the network journey? It's ambiguous.

As soon as you have to do something complicated the Fetch options could get messy, but we can set them when the Fetcher is created rather than when it makes the fetch.