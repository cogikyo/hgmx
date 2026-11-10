import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const body = el.querySelector('path[d^="m14 20"]');
  animate(body, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const bars = [el.querySelector('path[d="M6 4h12"]'), el.querySelector('path[d="M6 8h12"]')];
  animate(bars, { scaleX: [0, 1.25, 1] }, { duration: τmid, delay: stagger(0.06, { startDelay: τmin }) });
}
