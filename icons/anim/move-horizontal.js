import { animate, τmin } from "../motion.js";

export default function (el) {
  const nudge = { duration: 0.3 };
  animate(el.querySelector('path[d^="m18 8"]'), { translateX: [0, 3, -1, 0] }, { ...nudge, delay: 0 });
  animate(el.querySelector('path[d^="m6 8"]'), { translateX: [0, -3, 1, 0] }, { ...nudge, delay: τmin });
}
