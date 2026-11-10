import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const axes = el.querySelector('path[d^="M3 3"]');
  animate(axes, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const line = el.querySelector('path[d^="m19 9"]');
  animate(line, { scaleY: [0, 1.25, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
