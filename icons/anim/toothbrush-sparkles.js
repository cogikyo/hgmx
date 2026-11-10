import { animate, stagger, τmax, τmin } from "../motion.js";

export default function (el) {
  const brush = el.querySelector("g:not([transform-origin])");
  animate(brush, { translateX: [0, 1.5, -1.5, 0], translateY: [0, -1.35, 1.35, 0] }, { duration: τmax });

  const sparkles = el.querySelectorAll("g[transform-origin]");
  animate(sparkles, { scale: [0, 1.25, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
