import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const disk = el.querySelector('path[d^="M15.2"]');
  animate(disk, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const label = el.querySelector('path[d^="M17 21"]');
  animate(label, { scaleY: [0, 1.25, 1] }, { duration: τmid, delay: τmin });
}
