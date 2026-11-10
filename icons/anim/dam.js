import { animate, stagger } from "../motion.js";

export default function (el) {
  const waves = el.querySelectorAll('path[d^="M11 11.31"], path[d^="M11.75 18"]');
  animate(waves, { pathLength: [0, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.12) });
}
