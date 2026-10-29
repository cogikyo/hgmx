import { animate, stagger, τmid } from "../motion.js";

export default function (el) {
  const bubble = el.querySelector('path[d^="M2.992"]');
  animate(bubble, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const dots = el.querySelectorAll('path[d$="h.01"]');
  animate(dots, { translateY: [0, -2, 1, 0] }, { duration: 0.3, delay: stagger(0.06) });
}
