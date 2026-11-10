import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const lines = el.querySelectorAll('path[d$="H3"]');
  animate(lines, { translateX: [-3, 0], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.06) });

  const tick = el.querySelector('path[d^="m15 18"]');
  animate(tick, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
  animate(tick, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
