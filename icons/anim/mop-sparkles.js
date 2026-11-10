import { animate, stagger, τmax, τmin } from "../motion.js";

export default function (el) {
  const mop = el.querySelector("g.mop");
  animate(mop, { translateX: [0, -2.5, 1.5, 0] }, { duration: τmax });

  const sparkles = el.querySelectorAll("g:not(.mop)");
  animate(sparkles, { scale: [0, 1.25, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
