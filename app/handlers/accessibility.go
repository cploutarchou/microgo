package handlers

import (
	"context"
	"net/http"
)

func (h *Handlers) render(w http.ResponseWriter, r *http.Request, tmpl string, vars, data interface{}) error {
	return h.APP.Render.Page(w, r, tmpl, vars, data)
}

func (h *Handlers) sessionPut(ctx context.Context, key string, value interface{}) {
	h.APP.Session.Put(ctx, key, value)
}
func (h *Handlers) sessionExist(ctx context.Context, key string) bool {
	return h.APP.Session.Exists(ctx, key)
}
func (h *Handlers) sessionGet(ctx context.Context, key string) interface{} {
	return h.APP.Session.Get(ctx, key)
}

func (h *Handlers) sessionRemove(ctx context.Context, key string) {
	h.APP.Session.Remove(ctx, key)

}
func (h *Handlers) sessionRenew(ctx context.Context) error {
	return h.APP.Session.RenewToken(ctx)

}
func (h *Handlers) sessionDestroy(ctx context.Context) error {
	return h.APP.Session.Destroy(ctx)

}

func (h *Handlers) createRandomString(n int) string {
	return h.APP.CreateRandomString(n)
}
