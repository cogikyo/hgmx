import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const left = el.querySelector('path[d^="M9 17"]');
  animate(left, { translateX: [-1.5, 0.5, 0] }, { duration: τmid });
  const right = el.querySelector('path[d^="M15 7"]');
  animate(right, { translateX: [1.5, -0.5, 0] }, { duration: τmid });

  const slash = el.querySelector('line[x1="2"]');
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
