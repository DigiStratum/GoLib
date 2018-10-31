package resources

/*
Map of resource path to the Resource and its properties

Note: Our path follows the following general structure:

scope/language/resource_relative_path_and_filename

Such that:

scope = "public" (for things that a client may request directly), or "private" (internal use only)
language = "xxx_YY" where xxx=country code and YY=territory/locale (e.g. "en-US") or "default" for any mismatch
resource_relative_path_and_filename = as described, a customary, relative path and filename

We could potentially supply additional scopes other than public and private, but those would be on
the implementation to structure (ResourceManager.GetScopedResource() supports this). We could also
support any language identifier scheme, but we stick to xxx_YY (or just xx_YY) to standardize.

ref: https://tools.ietf.org/html/rfc5646
ref: http://cldr.unicode.org/

TODO: Add some supporting funcs to Resource to get a list of Resources below a given path (i.e. everything in a dir)

*/

type ResourceMap map[string]*Resource

func (rm ResourceMap) GetResource(path string) *Resource {
	if r, ok := rm[path]; ok {
		return r
	}
	return nil
}

