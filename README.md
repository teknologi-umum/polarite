<kbd>
<img src="./.github/logo.png" width="1000px">
</kbd>

<div style="text-align:center;">
<strong>Polarite</strong> is a Pastebin alternative made for simplicity written in Go.
</div>

## Usage

### Web Interface

Soon.

### API

Send a `POST` request to `link coming soon` with:

- `Content-Type` header with the value of `text/plain`
- Non-empty `request body` consisting of the text you want to store
- `Authorization` header with the value of `Your Name <your email>`

Example request:

- cURL
```sh
curl \
    -H "Content-Type: text/plain" \
    -H "Authorization: John Doe <john@example.com>" \
    -X POST \
    -d "Hello World" <link coming soon>
```

- Javascript (via Fetch API)
```js
fetch("link coming soon", {
    method: "POST",
    headers: {
        "Content-Type": "text/plain",
        "Authorization": "John Doe <john@example.com>"
    },
    body: "Hello world!"
})
```

- Go
```go
import (
    "net/http"
    "strings"
)

func Polarite() {
    body := strings.NewReader("Hello world")
    req, err := http.NewRequest(http.MethodPost, "link coming soon", body)
    req.Header.Add("Content-Type", "text/plain")
    req.Header.Add("Authorization", "John Doe <john@example.com>")

    client := &http.Client{}
    resp, err := client.Do(req)
}
```

- C#
```c#
using System.Net.Http;

var client = new HttpClient();
var request = new HttpRequestMessage() {
    RequestUri = new Uri("link coming soon"),
    Method = HttpMethod.Post,
    Headers = {
        { "Authorization", "John Doe <john@example.com>" },
        { "ContentType", "text/plain" }
    },
    Content = new StringContent("Hello world", Encoding.UTF8, "text/plain")
};

var task = await client.SendAsync(request);
```

## I'm here for Hacktoberfest, what can I do?

If you're new to open source, we really recommend reading a few articles about contributing to open source projects:

- [Open Source Guide's How to Contribute to Open Source](https://opensource.guide/how-to-contribute/)
- [Hacktoberfest Contributor's Guide: How To Find and Contribute to Open-Source Projects](https://www.digitalocean.com/community/tutorials/hacktoberfest-contributor-s-guide-how-to-find-and-contribute-to-open-source-projects)
- [Tips for high-quality Pull Request](https://twitter.com/sudo_navendu/status/1437456596473303042)

You can start by reading the [CONTRIBUTING.md](./CONTRIBUTING.md). Then you can search for [issues that you can work on](https://github.com/teknologi-umum/polarite/issues?q=is%3Aopen+is%3Aissue+label%3Ahacktoberfest).

Have fun!

## Why the name, Polarite?

In the dawn of time, it began with the birth of [Graphene](https://github.com/teknologi-umum/graphene) repository, which its' name was picked from the name of a mineral.
Then, not so long after, another repository called [Flourite](https://github.com/teknologi-umum/flourite) emerged. It's actually a typo of Fluorite, another name of a mineral.
Now, where mankind stands, we want to continue that convention, to pick a name from a [list of mineral on Wikipedia](https://en.wikipedia.org/wiki/List_of_minerals).
