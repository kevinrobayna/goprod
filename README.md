# GoProd

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/kevinrobayna/goprod)

Sample application written with Golang using web frameworks that would replace the ones I'm used to and love from the
JVM world.

Personally I love the following:

* [Dropwizard.io](https://www.dropwizard.io/en/latest/)
    * It's a really simple framework that does very few things but tries to guide you building a minimal but clean
      microservice.
* [Guice](https://github.com/google/guice)
    * Google says that it's a lightweight dependency injection framework, and that's true. It's very simple and enables
      you to do almost everything you would ever want, including things you should not do.
* [Hibernate](https://hibernate.org)/[Gorm](https://gorm.grails.org)
    * Both are ORM's. Gorm is a lightweight ORM that is very easy to use and has a very small footprint. Hibernate is a
      full-featured ORM that is very powerful and has a large footprint.
    * They are both excellent ORMs, and they don't tend to get in the way.
* [SLF4J (Simple Log Facade For Java)](http://www.slf4j.org)
    * Clean API than enables many libraries to log messages without any issues.
    * Go does not provide an interface for logs, instead it just provides classes. This means that you would
      implement whatever you need. I don't really like that as moving from one library to another might become a pain as
      you need to glue them together.

## Goals

With this app I want to build something using all the standards I've learned to like in my career. This includes
not only the above but building in a clean way.

This app should be observable at all times and be able to scale up in terms of features and teams working on this.

Additionally, everything should be tested with different layers of tests from unit to e2e and contract tests so that
others can rely on our API to be dependable.
