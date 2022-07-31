#Object

Objects are blocks of data that we manage whether from disk, a database, or some remote service.
We provide a standard interface to working with Objects as well as supporting libraries because,
in the cloud, Objects tend to be located in a whole variety of places - if we end up writing the
application in such a way as to expect a given Object to be accessed in a specific way, but later
want to redesign to access it a different way, then abstracting the implementation away from the
underlying method of access is the best way to future-proof our design.

In the spirit of Cloud Native RESTful thinking, an Object is identified with a "path" which is
representative of a directory hierarchy followed by a filename. Whether the path is a literal file
and directory structure on a disk, a virtual representation of such in a service, or simply a
unique key to a record in a database, the application will not know or care. But the essence of a
Object is that it is a "thing" that can be handled, moved around, used, read/written, and all the
rules of REST and HTTP response codes apply (i.e. if you ask for a Object which doesn't exist, you
get a 404 response).

An Object path follows the general structure:

SCOPE/CONTEXT/LANGUAGE/object_relative_path_and_filename

Such that:
* SCOPE = "public" (anonymous client), "private" (internal client), "protected" (authorized client)
* CONTEXT = Any contextualization path component between the scope and the language-specific
objects; (may be empty)
* LANGUAGE = "xxx-YY" where xxx=country code and YY=territory/locale (e.g. "en-US") or "default"
for any mismatch object_relative_path_and_filename = as described, a customary, relative path and
filename; the \-YY suffix is optional (e.g. "en" is acceptable)

FIXME: change language specifier to ISO/RFC standard.

We could potentially supply additional scopes other than public and private, but those would be on
the implementation to structure (ObjectManager.GetScopedObject() supports this). We could also
support any language identifier scheme, but we stick to xxx_YY (or just xx_YY, or xx) to standardize.

ref: https://tools.ietf.org/html/rfc5646
ref: http://cldr.unicode.org/

