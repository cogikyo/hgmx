import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const body = [el.querySelector('path[d^="M19 13"]'), el.querySelector('path[d="M5 20h14"]')];
  animate(body, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const bars = [el.querySelector('path[d="M10 10V3"]'), el.querySelector('path[d="M14 10V3"]')];
  animate(bars, { scaleY: [0, 1.25, 1] }, { duration: τmid, delay: stagger(0.06, { startDelay: τmin }) });
}
