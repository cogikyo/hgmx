import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("rect");
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const bar = el.querySelector('path[d^="M7 9"]');
  animate(bar, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const stem = el.querySelector('path[d="M12 7v10"]');
  animate(stem, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
  const foot = el.querySelector('path[d="M10 17h4"]');
  animate(foot, { scaleX: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
