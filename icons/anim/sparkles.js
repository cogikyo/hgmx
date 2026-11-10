import { animate, stagger, τmax, τmin } from "../motion.js";

export default function (el) {
  const spark = el.querySelector('path[d^="M11.017"]');
  animate(spark, { scale: [0.5, 1.25, 1], rotate: [0, 90] }, { duration: τmax });

  const glints = el.querySelectorAll('path:not([d^="M11.017"]), circle');
  animate(glints, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
