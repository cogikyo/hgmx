import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector('path[d^="M5 21a2 2 0 0 1-2-2V5"]');
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const dashes = [el.querySelector('path[d="M14 21h1"]'), el.querySelector('path[d="M9 21h1"]')];
  animate(dashes, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmid }) });

  const open = el.querySelector('path[d^="M10 9.5"]');
  animate(open, { translateX: [2, -1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const close = el.querySelector('path[d^="m14 9.5"]');
  animate(close, { translateX: [-2, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
