package dependencies

/*
DependencyInjectorIfc is meant to indicate that a given implementation can inject itself into a
DependencyInjectableIfc that is prepared to receive it. This moves a portion of the DI model into
each resource that is anticipated to be injected into various consumers.

TODO: Consider: is this A Bad Idea (tm)? This spreads DI concerns into the interface of the library
resource that may be better off consolidated into the parent logic setting up the DI: instead of
one place doing all the setup and injection, there are many. The good thing here is that it is an
optional construct for our DI model, it is not so dissimilar from requiring OOP classes to have
a constructor in simply becoming part of the design paradigm for object instance lifecycle
management, and there is potential here to syntactically streamline the implementation so that
DI logic is more readable, easier to grok. Let's see how it goes before we decide whether to
accept or reject into our overall DI scheme as a long-term resident.

*/

type DependencyInjectorIfc interface {
	InjectInto(client DependencyInjectableIfc, variant ...string) error
}

