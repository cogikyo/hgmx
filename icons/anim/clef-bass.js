import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const head = el.querySelector("circle");
  animate(head, { scale: [1, 1.3, 1] }, { duration: τmin });
  const curve = el.querySelector('path[d^="M5 8"]');
  animate(curve, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });
  const dots = [el.querySelector('path[d="M19 6h.01"]'), el.querySelector('path[d="M19 11h.01"]')];
  animate(dots, { scale: [1, 1.5, 1] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmid }) });
}
