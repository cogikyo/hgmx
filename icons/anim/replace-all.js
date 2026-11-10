import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const arrow = el.querySelectorAll('path[d^="m3 7"], path[d^="M6 10V5"]');
  animate(arrow, { translateY: [0, 2, -0.5, 0] }, { duration: τmid });
  const targets = [el.querySelector("rect"), el.querySelector('path[d^="M14 14"]'), el.querySelector('path[d^="M19 14"]')];
  animate(targets, { scale: [1, 0.85, 1.1, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
