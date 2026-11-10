import { animate, stagger, τmax, τmin } from "../motion.js";

export default function (el) {
  const [upper, lower] = el.querySelectorAll("g");
  animate(upper, { rotate: [0, 15, -5, 0] }, { duration: τmax });
  animate(lower, { rotate: [0, -15, 5, 0] }, { duration: τmax });

  const dashes = el.querySelectorAll('path[d$="12h-2"]');
  animate(dashes, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.12, { startDelay: τmin }) });
}
