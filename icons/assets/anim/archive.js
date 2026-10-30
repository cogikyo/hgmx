import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const lid = el.querySelector("rect");
  animate(lid, { translateY: [0, -3, 1, 0] }, { duration: τmax });

  const box = el.querySelector('path[d^="M4 8"]');
  animate(box, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const slot = el.querySelector('path[d="M10 12h4"]');
  animate(slot, { scaleX: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
