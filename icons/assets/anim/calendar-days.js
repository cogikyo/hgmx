import { animate, stagger, τmid } from "../motion.js";

export default function (el) {
  const page = el.querySelector("rect");
  animate(page, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const days = el.querySelectorAll('path[d$="h.01"]');
  animate(days, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.03) });
}
