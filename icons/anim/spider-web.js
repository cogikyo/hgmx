import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const outer = el.querySelector('path[d^="M21 17.2"]');
  animate(outer, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const inner = el.querySelector('path[d^="M17 14.9"]');
  animate(inner, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
