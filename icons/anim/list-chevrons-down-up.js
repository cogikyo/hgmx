import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const lines = el.querySelectorAll('path[d$="h8"]');
  animate(lines, { translateX: [-3, 0], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.06) });

  const top = el.querySelector('path[d^="m15 5 "]');
  animate(top, { translateY: [0, 1.5, -0.5, 0] }, { duration: τmid, delay: τmin });
  const bottom = el.querySelector('path[d^="m15 19 "]');
  animate(bottom, { translateY: [0, -1.5, 0.5, 0] }, { duration: τmid, delay: τmin });
}
