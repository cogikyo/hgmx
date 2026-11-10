import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const tub = el.querySelector('path[d^="M4 12"]');
  animate(tub, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const bubbles = [el.querySelector('circle[cx="16.5"]'), el.querySelector('circle[cx="11.5"]'), el.querySelector('path[d="M15 3h.01"]')];
  animate(bubbles, { translateY: [0, -2, 0.5, 0] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
