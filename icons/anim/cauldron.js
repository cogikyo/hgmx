import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const drips = el.querySelectorAll('path[d="M7 12V8"], path[d="M11 10V8"]');
  animate(drips, { scaleY: [0, 1] }, { duration: τmid });
  const bubbles = [el.querySelector('circle[r="2"]'), el.querySelector('circle[r=".5"]'), el.querySelector('path[d="M19 2h.01"]')];
  animate(bubbles, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
