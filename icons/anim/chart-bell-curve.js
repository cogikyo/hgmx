import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const axes = el.querySelector('path[d^="M3 3"]');
  animate(axes, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const curve = el.querySelector('path[d^="M7 16c"]');
  animate(curve, { scaleY: [0, 1.25, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
