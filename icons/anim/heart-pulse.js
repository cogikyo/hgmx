import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const pulse = el.querySelector('path[d^="M3.22"]');
  animate(pulse, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmax });

  const heart = el.querySelector('path[d^="M2 9.5"]');
  animate(heart, { scale: [1, 1.15, 0.95, 1] }, { duration: τmid, delay: τmin });
}
