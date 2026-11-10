import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const folder = el.querySelector('path[d^="M9 20"]');
  animate(folder, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const branch = el.querySelector('path[d^="M18 19"]');
  animate(branch, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const target = el.querySelector('circle[cx="20"]');
  animate(target, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
