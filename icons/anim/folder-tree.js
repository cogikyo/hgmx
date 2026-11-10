import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const branches = el.querySelectorAll('path[d^="M3 5"], path[d^="M3 3"]');
  animate(branches, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const folders = el.querySelectorAll('path[d^="M20 10"], path[d^="M20 21"]');
  animate(folders, { translateX: [-4, 1, 0], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
