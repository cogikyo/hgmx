import { animate, stagger, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const badge = el.querySelector('path[d^="M3.85"]');
  animate(badge, { scale: [1, 0.85, 1.1, 1] }, { duration: τmax });

  const glyph = el.querySelectorAll("line");
  animate(glyph, { translateY: [-4, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: stagger(τmin / 2, { startDelay: τmin }) });
}
