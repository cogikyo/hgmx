import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const page = el.querySelector("rect");
  animate(page, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const rings = el.querySelectorAll('path[d$="v3"]');
  animate(rings, { translateY: [-3, 1, 0], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmid }) });
  const days = el.querySelectorAll('path[d$="h.01"]');
  animate(days, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.03) });
}
