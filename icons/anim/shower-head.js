import { animate, stagger } from "../motion.js";

export default function (el) {
  const drops = el.querySelectorAll('path[d$="v.01"]');
  animate(drops, { translateX: [-2, 0.5, 0], translateY: [-2, 0.5, 0], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.03) });
}
