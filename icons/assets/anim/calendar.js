import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const page = el.querySelector("rect");
  animate(page, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const rings = el.querySelectorAll('path[d$="v3"]');
  animate(rings, { translateY: [-3, 1, 0], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmid }) });

  const rule = el.querySelector('path[d="M3 9h18"]');
  animate(rule, { scaleX: [0, 1] }, { duration: τmid, delay: τmin });
}
