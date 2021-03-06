package wfe

import (
	"github.com/lyraproj/issue/issue"
	"github.com/lyraproj/puppet-evaluator/eval"
	"github.com/lyraproj/puppet-evaluator/types"
	"github.com/lyraproj/servicesdk/serviceapi"
	"github.com/lyraproj/servicesdk/wfapi"
	"github.com/lyraproj/wfe/api"
	"github.com/lyraproj/wfe/service"
)

type stateHandler struct {
	Activity
	typ eval.ObjectType
}

func StateHandler(def serviceapi.Definition) api.Activity {
	a := &stateHandler{}
	a.Init(def)
	return a
}

func (a *stateHandler) Init(def serviceapi.Definition) {
	a.Activity.Init(def)
	// TODO: Type validation. The typ must be an ObjectType implementing read, upsert, and delete.
	a.typ = service.GetProperty(def, `interface`, types.NewTypeType(types.DefaultObjectType())).(eval.ObjectType)
}

func (a *stateHandler) Run(c eval.Context, input eval.OrderedMap) eval.OrderedMap {
	ac := service.ActivityContext(c)
	op := service.GetOperation(ac)
	invokable := a.GetService(c)

	switch op {
	case wfapi.Read:
		return invokable.Invoke(c, a.name, `read`, input).(eval.OrderedMap)

	case wfapi.Upsert:
		return invokable.Invoke(c, a.name, `upsert`, input).(eval.OrderedMap)

	case wfapi.Delete:
		invokable.Invoke(c, a.name, `delete`, input)
		return eval.EMPTY_MAP
	default:
		panic(eval.Error(wfapi.WF_ILLEGAL_OPERATION, issue.H{`operation`: op}))
	}
}

func (a *stateHandler) Label() string {
	return ActivityLabel(a)
}

func (a *stateHandler) Style() string {
	return `stateHandler`
}
