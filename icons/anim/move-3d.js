import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const depth = el.querySelector('path[d^="m5 19"]');
  animate(depth, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const nudge = { duration: 0.3 };
  animate(el.querySelector('path[d^="m2 6"]'), { translateY: [0, -3, 1, 0] }, { ...nudge, delay: 0 });
  animate(el.querySelector('path[d^="m18 16"]'), { translateX: [0, 3, -1, 0] }, { ...nudge, delay: τmin });
}
