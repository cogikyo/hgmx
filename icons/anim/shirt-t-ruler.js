import { animate, stagger, τmax, τmin } from "../motion.js";

export default function (el) {
  const shirt = el.querySelector("g");
  animate(shirt, { scaleY: [1, 1.12, 0.96, 1] }, { duration: τmax });

  const ticks = el.querySelectorAll('path[d="M10 16v2"], path[d="M14 16v2"]');
  animate(ticks, { scaleY: [0, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.12, { startDelay: τmin }) });
}
