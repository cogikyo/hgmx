import { animate, stagger, τmax, τmid } from "../motion.js";

export default function (el) {
  const folder = el.querySelector('path[d^="M20.9"]');
  animate(folder, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const teeth = el.querySelectorAll('path[d="M15 11v-1"], path[d="M15 17v-2"]');
  animate(teeth, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });

  const pull = el.querySelector("circle");
  animate(pull, { translateY: [-8, 1, 0] }, { duration: τmax });
}
