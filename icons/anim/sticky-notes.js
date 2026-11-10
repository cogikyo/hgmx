import { animate, stagger, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const stack = el.querySelector("g");
  animate(stack, { rotate: [0, -12, 4, 0], translateY: [0, -4, 1, 0] }, { duration: τmax });

  const folds = [el.querySelector('path[d^="M10 8v5"]'), el.querySelector('path[d^="M16 2v5"]')];
  animate(folds, { scale: [0, 1.25, 1] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmid }) });
}
