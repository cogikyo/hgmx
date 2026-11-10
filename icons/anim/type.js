import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const bar = el.querySelector('path[d^="M4 7"]');
  animate(bar, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const stem = el.querySelector('path[d="M12 4v16"]');
  animate(stem, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
  const foot = el.querySelector('path[d="M9 20h6"]');
  animate(foot, { scaleX: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
