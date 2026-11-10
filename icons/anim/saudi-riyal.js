import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const stems = el.querySelectorAll('path[d^="M10 2.977"], path[d^="M14 4"]');
  animate(stems, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const bars = [el.querySelector('path[d="M19 8 4 14"]'), el.querySelector('path[d="m19 18-5 2"]')];
  animate(bars, { scale: [0, 1.25, 1] }, { duration: τmid, delay: stagger(0.06, { startDelay: τmin }) });
}
