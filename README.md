# sa-todo-list-api

## Hosted

This is hosted at: https://todoapi.benaychh.io/v1

## Running the API

1. Install Dependencies: `go get -d ./...`
1. Compile the program: `go build -o main .`
1. Setup a database using the `init.sql` script
1. Create a `.env` file using `.env.example` as a template
1. Load the arguments into your session using `export $(cat .env | xargs)`
1. Run the program: `./main`
1. ???
1. Profit

> Note: You can just use the `.env.example` file and setup your database that way if that ie easier, then run `export $(cat .env.example | xargs)`

You can also go the much easier way and just run `docker-compose build && docker-compose up` assuming that you have [docker-compose](https://docs.docker.com/compose/) installed. This is my preferred method.

> Note: For ease in testing, this is setup without a volume mount so the database is pretty ephemeral

## Running the tests

**Assumes you have already compiled the program**

1. Setup a database to match the `test.env` file
1. Load those variables: `export $(cat test.env | xargs)`
1. Run `go test -v`

> Note this will clear whatever table you are using

You can also go the much easier way and run the docker-compose file for the tests with the following command: `docker-compose -f run-tests.yml up --abort-on-container-exit --build`

## Thought Process

First of all, I want to be up-front about having never written an API in GO before and that I only started using the language on December 1st for Advent of Code. All of the code is my own (I did not ask any live person for help) but I did borrow heavily from the following:

1. https://itnext.io/structuring-a-production-grade-rest-api-in-golang-c0229b3feedc
1. https://itnext.io/building-restful-web-api-service-using-golang-chi-mysql-d85f427dee54
1. https://semaphoreci.com/community/tutorials/building-and-testing-a-rest-api-in-go-with-gorilla-mux-and-postgresql
1. http://go-database-sql.org
1. https://medium.com/statuscode/how-i-write-go-http-services-after-seven-years-37c208122831
1. http://blog.setale.me/2018/01/06/Go-Tests-multistage-docker-builds-pipeline/
1. https://www.cloudreach.com/blog/containerize-this-golang-dockerfiles/

After seeing a couple of different ways to implement an API, this seemed the most idomatic and felt like the code flowed pretty naturally from the setup.

I chose to use [Chi](https://github.com/go-chi/chi) because it was the first router I came across that had all of the features I needed and it reminded me of using express/hapi.

When it came to testing, I tried to take a mostly TDD approach but because I was learning as I was going, there were some times where I didn't know what I was testing until I had written a little bit of code. I think this is reflected in my git commit history.

On the topic of testing, I always tend to prefer Integration/E2E testing over Unit/Functional testing. In my opinion, Integration/E2E do the best job of protecting the user of the program while unit tests do a great job of protecting the code from the developer ;)

In this case, I chose to do integration tests with the database as I think it gives me the biggest bang for my buck in the short time I was working on this project. This gives me pretty good coverage and I was able to test all of my green paths and some of the most common red paths pretty easily. If I had more time, I would probably add some unit tests in there to get my coverage above 80% but testing is great until it becomes pedantic and I didn't want to go down a rabbit hole.

I ended up implementing 4 routes, getAll, createOne, updateOne and deleteOne. I considered using separate routes for completing vs updating but that didn't seem necessary. In retrospect, I am the least happiest with my updateOne route `(PATCHing to /${todoID})` as it felt like it was heading towards getting out of my control. With more time, I think I would probably ask for advice on that one and look at cleaning it up.

I would also like to bring in some documentation. When working with Nodejs, I can use the [Joi](https://github.com/hapijs/joi) package and then get swagger docs for free. I am sure there is something like that in Go as well but it seemed a bit out of scope for this - especially since I was the only one consuming the API and there are only 4 basic routes.

Finally, I chose to use Environment Variables for all forms of config so that the resulting docker image could be build once, deploy anywhere.

## More Time Wishlist

1. Coverage up to ~80%, I know it is an arbitrary number but it has historically served me well
1. Get some unit testing around `sendJson` and `sendError`
1. Payload validation library (incoming and outgoing) like Joi
1. Swagger (or the like) docs
1. Pair with an experienced Go dev for help on my patch route.
