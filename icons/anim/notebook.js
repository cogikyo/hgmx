import { animate, stagger, τmax } from "../motion.js";

export default function (el) {
  const cover = el.querySelectorAll('rect, path[d^="M16 2"]');
  animate(cover, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmax });

  const rings = el.querySelectorAll('path[d$="h4"]');
  animate(rings, { scaleX: [0, 1.25, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
