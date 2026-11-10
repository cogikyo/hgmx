import { animate, stagger, τmin } from "../motion.js";

export default function (el) {
  const one = el.querySelector("g");
  const two = el.querySelector('path[d^="M6.5"]');
  animate([one, two], { scale: [0, 1.25, 1] }, { duration: 0.3, delay: stagger(0.21) });

  const lines = el.querySelectorAll('path[d$="h10"]');
  animate(lines, { opacity: [0, 1], x: [-3, 0] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmin }) });
}
