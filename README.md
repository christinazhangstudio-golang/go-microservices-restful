# go-microservices-restful
If using for the first time, be sure to run 

```
npm install
```
and then 

```
yarn start
```

Docs on React here: https://github.com/facebook/create-react-app

The react app is running on localhost:3000, but the products-api is on localhost:9090/products.

When running the react app at first, the below error is presented:

```
"Access to XMLHttpRequest at 'http://localhost:9090/products' from origin 'http://localhost:3000' has been blocked by CORS policy: No 'Access-Control-Allow-Origin' header is present on the requested resource."
```

CORS blocks this, since it blocks requests from other origins besides the one its loaded. ReactJS is running on 3000, so origin is localhost:3000. A different port (i.e. 9090) is a different origin. CORS protects the browser gettiing requests from other origins (forwarding cookies causes security compromises). To solve this, CORS in gorilla/mux is used (see main.go).

Handling and serving files is a bit complex... refer to https://www.youtube.com/watch?v=ctmhYJpGsgU&list=PLmD8u-IFdreyh6EUfevBcbiuCKzFk0EW_&index=10
