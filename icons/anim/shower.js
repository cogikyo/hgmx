import { animate, stagger } from "../motion.js";

export default function (el) {
  const drops = el.querySelectorAll('path[d$="h.01"]');
  animate(drops, { translateY: [-3, 1, 0], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.03) });
}
