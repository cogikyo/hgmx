import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const badge = el.querySelector('path[d^="M3.85"]');
  animate(badge, { pathLength: [0, 1], opacity: [0.5, 1], rotate: [-180, 0] }, { duration: τmid });
  const glyph = el.querySelectorAll("line");
  animate(glyph, { translateY: [-4, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: stagger(τmin / 2, { startDelay: τmin }) });
}
