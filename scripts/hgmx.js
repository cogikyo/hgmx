{
  let api, latest, seq = 0;
  const navigation = ({ headers }) => headers["HX-Boosted"] || headers["HX-History-Restore-Request"];

  htmx.registerExtension("stale", {
    htmx_config_request: (elt, { ctx }) => {
      if (navigation(ctx.request)) ctx.seq = ++seq;
    },
    htmx_before_request: (elt, { ctx }) => {
      if (navigation(ctx.request) && ctx.seq > (latest?.seq ?? 0)) latest = ctx;
    },
    htmx_after_request: (elt, { ctx }) =>
      elt.isConnected && !ctx.request.signal?.aborted && (!navigation(ctx.request) || ctx === latest),
  });

  htmx.registerExtension("navigation", {
    init: (internal) => (api = internal),
    htmx_config_request: (elt, { ctx }) => {
      let { headers } = ctx.request;
      if (headers["HX-History-Restore-Request"]) {
        Object.assign(headers, api.HCON.parse(api.attributeValue(elt, "hx-config")).headers);
      }
    },
    htmx_after_swap: (elt, { ctx }) => {
      if (!navigation(ctx.request)) return;
      let main = document.querySelector("main");
      if (!main || main.contains(document.activeElement)) return;
      let target = main.querySelector("h1") ?? main;
      if (!target.hasAttribute("tabindex")) target.tabIndex = -1;
      target.focus({ preventScroll: true });
    },
  });
}
