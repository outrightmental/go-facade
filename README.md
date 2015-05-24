# Facade 

[![Build Status](https://travis-ci.org/outrightmental/facade.svg?branch=master)](https://travis-ci.org/outrightmental/go-facade)

Facade is as minimal as it gets: in essence, it parses a single HTML file, caches it in memory, and on demand will return that HTML file with one string of HTML content injected into one HTML element in the original page.

Facade provides a convenient way to serve a static frontend UI from a backend server.

Use any frontend that builds to a static build output folder for distribution (e.g. [Gulp-Bower-Angular](https://github.com/Swiip/generator-gulp-angular))

Use any webserver that allows different URLs to be routed to different files/servers. (e.g. [Nginx](http://nginx.org/))
    
### Why?

Separation of concerns, according to [The Unix Philosophy](http://en.wikipedia.org/wiki/Unix_philosophy.

Use [Nginx](http://nginx.org/) for what it's great at.

Use a static front end build workflow, and package management (e.g. [Bower](http://bower.io/)) for what it's great at.

And use Go to drive your application's logical backend without being concerned by ^

It's power is in *how* the entire stack is configured to use a webserver (e.g. [Nginx](http://nginx.org/)) as the first line of HTTP service, to take as much custom work as possible off the production servers.

Check out [this blog post](http://www.outrightmental.com/facade-painless-middleware-frontending-for-go/)
for explanation how Facade is different from other frontending solutions.

### Usage

Our user requests a particular URL, and the webserver (e.g. [Nginx](http://nginx.org/)):

  1. Searches in the frontend's static build output folder for distribution, to see if the requested URL matches an exact file. For example, `/` (matches `index.html`) or `/scripts/app-ddf26c4d.js` or `/styles/vendor-f18be9e6.css` or `/assets/images/logo.png`. If a match is found, the HTTP request is terminated here, by serving the static file.
  2. If no static file is found, pass the request through to the Go webservice ***with the URL completely intact***.

That's where Facade comes in. It will serve our same `index.html` *for every URL that it encounters*, with one twist. Inside of our `index.html` is an HTML element with a special attribute, `facade`. For example:

    <!doctype html>
    <html>
      <head>
        <title>Freshest</title>
        <link rel="stylesheet" href="/styles/vendor-f18be9e6.css">
        <link rel="stylesheet" href="/styles/app-25ab4bd1.css">
      </head>
      <body ng-app="freshest" ng-controller="MainCtrl">
        <div facade ui-view=""><!-- 

            **go-facade injects content into here!**

        --></div>
        <script src="/scripts/vendor-23555feb.js"></script>
        <script src="/scripts/app-e39a84b8.js"></script>
      </body>
    </html>

So when our user visits

    https://freshest.io/oauth/authorize

Nginx passes the request through to our Go webserver- perhaps [Gorilla Mux](http://www.gorillatoolkit.org/pkg/mux) -actions are performed silently, followed by the construction of our Go service's HTTP Response.

    distPath := Getenv('PATH_TO_FRONTEND_INDEX_DOT_HTML')
    frontend := facade.New(distPath)

Our frontend has been built via any method, and Facade is configured to see only its final built distribution.

And when it's time to provide the HTTP response:

    content := "<p>This will be injected into the HTML element in our page with the attribute `facade`</p>"
    responsewriter.Header().Set( "Content-Type", "text/html")
    responsewriter.WriteHeader( result.Code)
    fmt.Fprint( responsewriter, frontend.Write(content))

### Development

**Next Up: Also inject header tags, e.g. `<title>` and `<description>`**

### Contributing

0. Find an issue that bugs you / open a new one.
1. Discuss.
2. Branch off, commit, test.
3. Make a pull request / attach the commits to the issue.
