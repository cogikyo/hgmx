import { animate, stagger, τmax, τmin } from "../motion.js";

export default function (el) {
  const map = el.querySelector("g");
  animate(map, { scaleX: [0.25, 1.15, 1] }, { duration: τmax });

  const folds = [el.querySelector('path[d^="M9 3"]'), el.querySelector('path[d^="M15 5"]')];
  animate(folds, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmin }) });
}
