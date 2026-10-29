const reduce = matchMedia("(prefers-reduced-motion: reduce)");

document.addEventListener("hovered", async (e) => {
  const svg = e.target.closest?.("svg[data-icon]");
  if (!svg || reduce.matches) return;
  const { default: play } = await import(`./anim/${svg.dataset.icon}.js`);
  play(svg);
});
