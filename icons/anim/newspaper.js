import { animate, stagger, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const page = el.querySelector('path[d^="M4 22"]');
  animate(page, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmax });

  const headline = el.querySelector("rect");
  animate(headline, { scaleX: [0, 1.1, 1] }, { duration: τmid, delay: τmin });

  const lines = [el.querySelector('path[d^="M18 14"]'), el.querySelector('path[d^="M15 18"]')];
  animate(lines, { opacity: [0, 1], x: [-3, 0] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmid }) });
}
