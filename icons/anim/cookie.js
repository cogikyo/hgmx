import { animate, stagger, τmid } from "../motion.js";

export default function (el) {
  const cookie = el.querySelector('path[d^="M11.496"]');
  animate(cookie, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const chips = el.querySelectorAll('path[d$="h.01"]');
  animate(chips, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.03) });
}
