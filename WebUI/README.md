UI modeling; we declare/configure a bunch of resources from the top, down, and then, per request, use a bottom, up approach to selecting which of those resources to use to render the output for the resulting page (for requests that result in a page rendering):

TOP-DOWN
1. Reusable full page layout, styling, and/or code templates ("Layouts") that are broadly applicable
2. Reusable partial page layout, styling, and/or code templates ("Fragments") included into Layouts
3. Reusable complete or partial styling specifications ("Styles") specifying fonts, colors, graphics, layout/positioning details
4. Reusable, static components ("Widgets") with predefined behavior included into Layouts and Fragments from a library of selections, global or subscription based
5. Reusable presentation packaging ("Schemes") specifying Layout, Fragment(s), and Styles
6. Customization of named Schemes with preview and versioning
7. Expandable library of Layouts, Fragments, Styles, and Schemes which new sites clone to start, thumbnail representations; "starters", possibly packaged, free and paid

BOTTOM-UP
8. Module-selectable Scheme
9. URI mapping to module
10. Inclusion of layout, styling, code, and/or configuraiton values ("Module Content") into Scheme based on URI/request/payload
11. Inculsion of layout, styling, code, and/or configuration values ("Widget Content") based on Scheme's client-side dynamic rendering, API interaction

PAGE RENDERING
For a given client request which is going to deliver rendered page content as a response, we render as follows:
1. Domain map identifies Module based on URI, or "default"
2. Module Scheme identifies Layout, Fragment(s), and Style(s)
3. Read-through cache for Rendered Schemes renders and caches Scheme result as needed and delivers rendered document to Module
4. Module hydrates rendered document with Module Content via magic tag substitution
5. Hydrated document is returned to client
6. Client executes dynamic code as supplied to render additional, dynamic Widget Content

MAGIC TAGS
Stratify SaaS Platform supports recursive tag substitution into in-memory bodies of text. Currently this is utilized in configuration data, config.json, to resolve dynamic configuration elements which must propagate in order to avoid declaration redundancy, and thereby reducing potential for error by way of introducing inconsistencies. Substitution is performed by way of generating a name/value pair dictionary and repeatedly scanning the target document for occurrences of "Magic Tags", denoted as %NAME% within the document, such that NAME is the name of one of the dictionary entries, otherwise ignored. This process is repeated until all resolvable Magic Tags are resolved, or infinite recursion is detected. config.json is a bit of a special case in that the document is also the dictionary, so it is self-referential to update configuration values which are dependent upon other configuration values.

Dictionary name tags meant for direct substitution are the simplest form of Magic Tag. "Action Tags" are a specialized form of Magic Tag which contain specifiers that cause more sophisticated behavior than simple dictionary substitution to occur. Action Tags are denoted by separating an action to be performed from the tag name with a ":" such as:

"%frag:fragmentname%"

This would direct the document processor to include the contents of the named template fragment in this place within the document. Note that the document processor would need to be provided with a contextualized set of resources in order to have the named fragment available for inclusion.

