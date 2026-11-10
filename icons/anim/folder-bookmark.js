import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const folder = el.querySelector('path[d^="M20 20"]');
  animate(folder, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const ribbon = el.querySelector('path[d^="M12 6"]');
  animate(ribbon, { scaleY: [1, 1.15, 0.95, 1] }, { duration: τmid, delay: τmin });
}
