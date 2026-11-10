import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const bars = [el.querySelector('path[d="M6 4v16"]'), el.querySelector('path[d="M10 4v16"]')];
  animate(bars, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.06) });
  const brace = el.querySelector('path[d^="M14 4.764"]');
  animate(brace, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
