import { animate, stagger, τmax, τmin } from "../motion.js";

export default function (el) {
  const pencil = el.querySelector("g:not([transform-origin])");
  animate(pencil, { rotate: [-15, 15, -10, 0], translateX: [0, -3, 0] }, { duration: τmax });

  const sparkles = el.querySelectorAll("g[transform-origin]");
  animate(sparkles, { scale: [0, 1.25, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
