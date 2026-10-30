import { animate, stagger, τmin } from "../motion.js";

export default function (el) {
  const bullets = el.querySelectorAll('path[d$="h.01"]');
  animate(bullets, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });
  const lines = el.querySelectorAll('path[d$="h13"]');
  animate(lines, { translateX: [-3, 0], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmin }) });
}
