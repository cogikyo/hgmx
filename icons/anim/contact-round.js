import { animate, stagger, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const card = el.querySelector("rect");
  animate(card, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmax });

  const rings = [el.querySelector('path[d^="M8 2"]'), el.querySelector('path[d^="M16 2"]')];
  animate(rings, { translateY: [-3, 1, 0], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmin }) });

  const body = el.querySelector('path[d^="M17.915"]');
  animate(body, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const head = el.querySelector("circle");
  animate(head, { translateY: [1, -2, 0] }, { duration: τmin, delay: τmid });
}
