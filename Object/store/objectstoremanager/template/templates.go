package template

/*

Mustache template handling via Managed ObjectStore

FIXME:
 * Make sure that this is thread-safe for managing the cache structure

TODO:
 * Any reason for this usage of object store manager to actually live as part of the objects library? Should be part of
   some web/HTML/render library instead, perhaps?

*/

import(
	"fmt"
	"bytes"
	"errors"

	osm "github.com/DigiStratum/GoLib/Object/store/objectstoremanager"
	"github.com/DigiStratum/Go-cbroglie-mustache"
)

type TemplateIfc interface {
	HydrateTemplate(templateName string, language string, data *map[string]string) (*string, error)
}

type templateCache		map[string]*mustache.Template

type Template struct {
	cache			templateCache
	objectStoreManager	*osm.ObjectStoreManager
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

// Make a new one of these
func NewTemplates(objectStoreManager *osm.ObjectStoreManager) *Template {
	t := Template{
		cache:			make(templateCache),
		objectStoreManager:	objectStoreManager,
	}
	return &t
}

// -------------------------------------------------------------------------------------------------
// TemplatesIfc Public Interface
// -------------------------------------------------------------------------------------------------

// Hydrate a named mustache template with the supplied data (name-value pair strings as a map)
func (r *Template) HydrateTemplate(templateName string, language string, data *map[string]string) (*string, error) {
	var template, err = r.getTemplate(templateName, language)
	if nil != err { return nil, err }
	var renderedTemplate bytes.Buffer
	template.FRender(&renderedTemplate, data)
	rendered := renderedTemplate.String()
	return &rendered, nil
}

// -------------------------------------------------------------------------------------------------
// Template Private Implementation
// -------------------------------------------------------------------------------------------------

// Read-through cache of named mustache templates
func (r *Template) getTemplate(templateName string, language string) (*mustache.Template, error) {
	templateKey := language + "." + templateName

	// If it's already in the cache, just return it!
	if cachedTemplate, ok := r.cache[templateKey]; ok { return cachedTemplate, nil }

	// Resolve template Object
	object := r.objectStoreManager.FindTemplate(language, templateName)
	if nil == object { return nil, fmt.Errorf("Template (%s) not in object tree", templateName) }

	// Parse it
	templateContent := object.GetContent()
	template, err := mustache.ParseString(*templateContent)
	if nil == err { r.cache[templateKey] = template }
	return template, err
}
