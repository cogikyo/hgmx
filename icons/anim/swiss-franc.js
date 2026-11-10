import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const body = el.querySelector('path[d^="M10 20"]');
  animate(body, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const bars = [el.querySelector('path[d="M10 10h7"]'), el.querySelector('path[d="M6 16h9"]')];
  animate(bars, { scaleX: [0, 1.25, 1] }, { duration: τmid, delay: stagger(0.06, { startDelay: τmin }) });
}
