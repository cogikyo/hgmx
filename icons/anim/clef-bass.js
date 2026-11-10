import { animate, stagger, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const clef = el.querySelector("g");
  animate(clef, { rotate: [0, -12, 8, 0] }, { duration: τmax });

  const dots = [el.querySelector('path[d="M19 6h.01"]'), el.querySelector('path[d="M19 11h.01"]')];
  animate(dots, { scale: [1, 1.5, 1] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmid }) });
}
