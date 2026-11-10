import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const c = el.querySelector('path[d^="M19 6"]');
  animate(c, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const bars = [el.querySelector('path[d="M4 10h12"]'), el.querySelector('path[d="M4 14h9"]')];
  animate(bars, { scaleX: [0, 1.25, 1] }, { duration: τmid, delay: stagger(0.06, { startDelay: τmin }) });
}
