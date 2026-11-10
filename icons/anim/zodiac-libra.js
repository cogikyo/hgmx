import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const beam = el.querySelector('path[d^="M3 16"]');
  animate(beam, { rotate: [0, -14, 10, -5, 0] }, { duration: τmax });

  const base = el.querySelector('path[d="M3 20h18"]');
  animate(base, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });
}
