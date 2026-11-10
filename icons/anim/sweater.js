import { animate, stagger, τmax, τmin } from "../motion.js";

export default function (el) {
  const sweater = el.querySelectorAll("path");
  animate(sweater, { translateY: [0, -2, 0.5, 0] }, { duration: τmax });

  const knit = el.querySelectorAll('path[d^="m6 10 "], path[d^="m6 16 "]');
  animate(knit, { pathLength: [0, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.12, { startDelay: τmin }) });
}
