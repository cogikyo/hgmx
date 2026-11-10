import { animate, stagger, τmax, τmin } from "../motion.js";

export default function (el) {
  const shirt = el.querySelector("g");
  animate(shirt, { scaleY: [1, 0.88, 1.04, 1] }, { duration: τmax });

  const buttons = el.querySelectorAll('path[d="M10 11h.01"], path[d="M10 15h.01"]');
  animate(buttons, { scale: [1, 1.75, 1] }, { duration: 0.3, delay: stagger(0.12, { startDelay: τmin }) });
}
