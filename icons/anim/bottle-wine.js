import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const bottle = el.querySelector('path[d^="M10 3"]');
  animate(bottle, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const label = el.querySelector('path[d^="M17 13h-4"]');
  animate(label, { scaleX: [0, 1.15, 1] }, { duration: τmid, delay: τmin });
}
