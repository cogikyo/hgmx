import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const ring = el.querySelector('path[d^="M17.41"]');
  animate(ring, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const leaf = el.querySelector('path[d^="M16 8q"]');
  animate(leaf, { scale: [0, 1.15, 1] }, { duration: τmid, delay: τmin });
}
