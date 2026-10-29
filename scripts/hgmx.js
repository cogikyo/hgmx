htmx.registerExtension("stale", {
  htmx_after_request: (elt) => elt.isConnected,
});

{
  let api;
  htmx.registerExtension("navigation", {
    init: (internal) => (api = internal),
    htmx_config_request: (elt, { ctx }) => {
      let { headers } = ctx.request;
      if (headers["HX-History-Restore-Request"]) {
        Object.assign(headers, api.HCON.parse(api.attributeValue(elt, "hx-config")).headers);
      }
    },
    htmx_after_swap: (elt, { ctx }) => {
      let { headers } = ctx.request;
      if (!headers["HX-Boosted"] && !headers["HX-History-Restore-Request"]) return;
      let main = document.querySelector("main");
      if (!main || main.contains(document.activeElement)) return;
      let target = main.querySelector("h1") ?? main;
      if (!target.hasAttribute("tabindex")) target.tabIndex = -1;
      target.focus({ preventScroll: true });
    },
  });
}
