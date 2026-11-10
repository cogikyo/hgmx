import { animate, stagger, τmax, τmin } from "../motion.js";

export default function (el) {
  const wand = el.querySelector("g.wand");
  animate(wand, { rotate: [0, -10, 4, 0] }, { duration: τmax });

  const sparkles = el.querySelectorAll("g:not(.wand)");
  animate(sparkles, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
