import { animate, stagger, τmax, τmin } from "../motion.js";

export default function (el) {
  const broom = el.querySelector("g[transform-origin]");
  animate(broom, { rotate: [0, -10, 8, -4, 0] }, { duration: τmax });

  const sparkles = el.querySelectorAll("g:not([transform-origin])");
  animate(sparkles, { scale: [0, 1.25, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
