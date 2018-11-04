package resources

/*

Mustache template handling

FIXME: Make sure that this is thread-safe for managing the cache structure

*/

import(
	"fmt"
	"bytes"
	"errors"

	"github.com/DigiStratum/Go-cbroglie-mustache"
)

type templateCache	map[string]*mustache.Template

type TemplateData	map[string]string

type Templates struct {
	cache		templateCache
	repositoryManager	*RepositoryManager
}

// Make a new one of these
func NewTemplates(rm *RepositoryManager) *Templates {
	t := Templates{
		cache:			make(templateCache),
		repositoryManager:	rm,
	}
	return &t
}

// Hydrate a named mustache template with the supplied data
func (tpl *Templates) HydrateTemplate(templateName string, language string, data *TemplateData) (*string, error) {
	var template, err = tpl.getCachedTemplate(templateName, language)
	if nil != err { return nil, err }
	var renderedTemplate bytes.Buffer
	template.FRender(&renderedTemplate, data)
	rendered := renderedTemplate.String()
	return &rendered, nil
}

// Provide read-through cache of named mustache templates
func (tpl *Templates) getCachedTemplate(templateName string, language string) (*mustache.Template, error) {

	templateKey := language + "." + templateName

	// If it's already in the cache, just return it!
	cachedTemplate, ok := tpl.cache[templateKey]
	if ok {
		return cachedTemplate, nil
	}

	// Resolve template Resource
	resource := tpl.repositoryManager.GetTemplate(templateName, language)
	if nil == resource {
		return nil, errors.New(fmt.Sprintf("Template (%s) not in resource tree", templateName))
	}

	// Parse it
	templateContent := resource.GetContent()
	template, err := mustache.ParseString(*templateContent)
	if nil == err { tpl.cache[templateKey] = template }
	return template, err
}

