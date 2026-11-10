import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const lines = el.querySelectorAll('path[d$="h8"]');
  animate(lines, { translateX: [-3, 0], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.06) });

  const up = el.querySelector('path[d^="m15 8 "]');
  animate(up, { translateY: [0, -3, 1, 0] }, { duration: τmid, delay: τmin });
  const down = el.querySelector('path[d^="m15 16 "]');
  animate(down, { translateY: [0, 3, -1, 0] }, { duration: τmid, delay: τmin });
}
