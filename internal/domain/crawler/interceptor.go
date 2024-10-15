package crawler

import "net/url"

type Interceptor = func(htmlStr string, target *url.URL)

// TODO: Create a type that will multiplex the html STR in a channel with multiple workers
