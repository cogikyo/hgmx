import { animate, stagger, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const card = el.querySelector("rect");
  animate(card, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmax });

  const lines = el.querySelectorAll('path[d$="h2"]');
  animate(lines, { opacity: [0, 1], x: [-2, 0] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmin }) });

  const body = el.querySelector('path[d^="M13 19"]');
  animate(body, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const head = el.querySelector("circle");
  animate(head, { translateY: [1, -2, 0] }, { duration: τmin, delay: τmid });
}
