import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector('path[d^="M5 21a2 2 0 0 1-2-2V5"]');
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const dashes = [el.querySelector('path[d="M14 21h1"]'), el.querySelector('path[d="M9 21h1"]')];
  animate(dashes, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmid }) });
}
