package ui

import (
	"fmt"
	"io"
	"sort"

	"strings"

	"github.com/cppforlife/go-cli-ui/ui"
	uitable "github.com/cppforlife/go-cli-ui/ui/table"
	"github.com/fatih/color"
	krsc "github.com/nimakaviani/kapp/pkg/kapp/resources"
	"github.com/nimakaviani/knative-inspect/pkg/resources"
)

const (
	ServiceKind = "Service"
	ApiGroup    = "serving.knative.dev"
)

type TreeView struct {
	Source      string
	ResourceMap [][]krsc.Resource
	Sort        bool
	Verbose     bool
}

func (v TreeView) Print(ui ui.UI) {
	groupHeader := uitable.NewHeader("Group")
	groupHeader.Hidden = true

	versionHeader := uitable.NewHeader("Version")
	versionHeader.Hidden = !v.Verbose

	table := uitable.Table{
		Title:   fmt.Sprintf("Resources in %s", v.Source),
		Content: "resources",

		Header: []uitable.Header{
			groupHeader,
			uitable.NewHeader("Namespace"),
			uitable.NewHeader("Name"),
			versionHeader,
			uitable.NewHeader("Kind"),
			uitable.NewHeader("Ready"),
			uitable.NewHeader("Reason"),
		},

		FillFirstColumn: true,
	}

	for _, rm := range v.ResourceMap {
		rt := buildResourceTree(rm)
		v.addRows(&table, rt, 0, "", false)
		v.addBlankRow(&table)
	}

	ui.PrintTable(table)
}

type resourceTree interface {
	Resource() krsc.Resource
	Children() []resourceTree
}

type resourceTreeImpl struct {
	root     krsc.Resource
	children []resourceTree
}

func (r *resourceTreeImpl) Resource() krsc.Resource {
	return r.root
}

func (r *resourceTreeImpl) Children() []resourceTree {
	return r.children
}

type byUID []resourceTree

func (a byUID) Len() int           { return len(a) }
func (a byUID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byUID) Less(i, j int) bool { return a[i].Resource().UID() < a[j].Resource().UID() }

type resourceRef struct {
	self   krsc.Resource
	parent krsc.Resource
}

func buildResourceTree(resources []krsc.Resource) resourceTree {
	resourceMap := make(map[string]*resourceRef)
	for _, r := range resources {
		if r == nil {
			continue
		}

		resourceMap[r.UID()] = &resourceRef{self: r}
	}

	for _, r := range resources {
		if r == nil {
			continue
		}

		for _, owner := range r.OwnerRefs() {
			if p, ok := resourceMap[string(owner.UID)]; ok {
				resourceMap[r.UID()].parent = p.self
			}
		}
	}

	root := findRoot(resourceMap)
	return buildTree(root, resourceMap)
}

func findRoot(nodes map[string]*resourceRef) *resourceRef {
	for _, n := range nodes {
		if n.parent == nil && n.self.Kind() == ServiceKind && strings.Contains(n.self.APIGroup(), ApiGroup) {
			return n
		}
	}
	panic("no root found")
}

func findKids(r krsc.Resource, resourceMap map[string]*resourceRef) []*resourceRef {
	var kids []*resourceRef
	for _, n := range resourceMap {
		if n.parent == r {
			kids = append(kids, n)
		}
	}
	return kids
}

func buildTree(r *resourceRef, resourceMap map[string]*resourceRef) resourceTree {
	rt := resourceTreeImpl{
		root:     r.self,
		children: []resourceTree{},
	}

	kids := findKids(r.self, resourceMap)
	for _, k := range kids {
		childTree := buildTree(k, resourceMap)
		rt.children = append(rt.children, childTree)
	}

	return &rt
}

func (v TreeView) addBlankRow(table *uitable.Table) {
	row := []uitable.Value{
		uitable.NewValueString(" "),
		uitable.NewValueString(" "),
		uitable.NewValueString(" "),
		uitable.NewValueString(" "),
		uitable.NewValueString(" "),
		uitable.NewValueString(" "),
	}
	table.Rows = append(table.Rows, row)
}

func (v TreeView) addRows(table *uitable.Table, rt resourceTree, depth int, connector string, hasSibling bool) {
	resource := rt.Resource()
	cond := resources.Conditions{Resource: resource}

	msg := "False"
	ready, msg := cond.IsSelectedTrue([]string{"Ready"})
	if ready {
		msg = "True"
	}

	reason := cond.Reason("Ready")

	if hasSibling {
		connector += " |"
	} else {
		connector += " "
	}

	var delim string
	if depth > 0 {
		delim = fmt.Sprintf("%sL", connector)
		if !hasSibling {
			delim += "_"
		}
		delim += " "
	}

	row := []uitable.Value{
		uitable.NewValueString(""),
		uitable.NewValueString(resource.Namespace()),
		uitable.NewValueString(resource.Name()),
		uitable.NewValueString(resource.APIVersion()),
		ValueEncoded{
			S: delim + resource.Kind(),
			Func: func(str string, opts ...interface{}) string {
				result := fmt.Sprintf(str, opts...)
				return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(result, "|L", "├─"), "|", "│"), "L_", "└─")
			},
		},
		ValueEncoded{
			S: msg,
			Func: func(str string, opts ...interface{}) string {
				result := fmt.Sprintf(str, opts...)
				if ready {
					return color.New(color.FgGreen).Sprintf("%s", result)
				}
				if msg != "" {
					return color.New(color.FgRed).Sprintf("%s", result)
				}
				return result
			},
		},
		ValueEncoded{
			S: reason,
			Func: func(str string, opts ...interface{}) string {
				result := fmt.Sprintf(str, opts...)
				if !ready && msg != "" {
					return color.New(color.FgRed).Sprintf("%s", result)
				}
				return ""
			},
		},
	}

	table.Rows = append(table.Rows, row)

	sortedChildren := byUID(rt.Children())
	sort.Sort(sortedChildren)
	for i, c := range sortedChildren {
		sibling := false
		if i+1 < len(rt.Children()) {
			sibling = true
		}
		v.addRows(table, c, depth+1, connector, sibling)
	}
}

type ValueEncoded struct {
	S    string
	Func func(string, ...interface{}) string
}

func (t ValueEncoded) String() string                  { return t.S }
func (t ValueEncoded) Value() uitable.Value            { return t }
func (t ValueEncoded) Compare(other uitable.Value) int { panic("Never called") }

func (t ValueEncoded) Fprintf(w io.Writer, pattern string, rest ...interface{}) (int, error) {
	return fmt.Fprintf(w, "%s", t.Func(pattern, rest...))
}
