import { animate, stagger, τmin } from "../motion.js";

export default function (el) {
  const lines = el.querySelectorAll('path[d$="H3"]');
  animate(lines, { scaleX: [0, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });
  const pages = el.querySelectorAll('path[d$="h.01"]');
  animate(pages, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
