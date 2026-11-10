import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const shaft = el.querySelector('path[d^="M22 18"]');
  animate(shaft, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const heads = el.querySelectorAll('path[d$="5-5"]');
  animate(heads, { translateX: [4, -2, 0], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
